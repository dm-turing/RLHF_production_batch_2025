package handlers

import (
	"Week_3/512415/turn2modela/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var uploadedfiles map[int]string

func ConvertFile(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get the file from the form
	file, handler, err := r.FormFile("uploadedfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	tempFilePath := filepath.Join("uploads", handler.Filename)
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the contents from the source to the destination
	src, _ := handler.Open()
	dst, _ := os.Create(tempFilePath)
	if _, err := io.Copy(dst, src); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	source := handler.Filename
	destination := strings.Replace(source, filepath.Ext(source), ".png", 1)
	conversion := &models.Conversion{
		UserID:      1, // hardcoded for demo purposes
		Source:      source,
		Destination: destination,
		Status:      "In progress",
	}
	uploadedfiles = map[int]string{}
	uploadedfiles[int(conversion.UserID)] = source
	StoreConversion(conversion)

	var id int64 // Modify this to a channel for concurrent access
	if id = convert(tempFilePath, destination); id == 0 {
		updateConversionStatus(conversion.ID, "Failed")
		log.Println("Conversion failed:", err)
	} else {
		updateConversionStatus(conversion.ID, "Completed")
	}

	if id != 0 {
		w.WriteHeader(http.StatusAccepted)
		downloadURL := fmt.Sprintf("/download/%d", id)
		response := map[string]string{"status": "accepted", "download_url": downloadURL}

		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "file conversion failed", http.StatusInternalServerError)
	}
}

func convert(source, destination string) int64 {
	// For simplicity, this function will only convert JPEG to PNG.
	// In a real application, you'd use a library like image/jpeg, image/png, etc.
	if strings.HasSuffix(source, ".jpg") {
		os.Rename(source, destination)
		return getID()
	}

	return 0 // Return 0 to indicate a conversion failure
}

func updateConversionStatus(id int64, status string) {
	db.Model(&models.Conversion{ID: id}).Update("status", status)
}

func getID() int64 {
	return 1 // For simplicity, always return a fixed ID in testing
}
