package handlers

import (
	"Week_3/512415/turn1modela/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ConvertFile(w http.ResponseWriter, r *http.Request) {
	if !Authenticate(w, r) {
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
		Source:      source,
		Destination: destination,
		Status:      "In progress",
	}

	StoreConversion(conversion)

	go func() {
		if err := convertImage(tempFilePath, destination); err != nil {
			updateConversionStatus(conversion.ID, "Failed")
			log.Println("Conversion failed:", err)
		} else {
			updateConversionStatus(conversion.ID, "Completed")
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

func convertImage(source, destination string) error {
	// For simplicity, this function will only convert JPEG to PNG.
	// In a real application, you'd use a library like image/jpeg, image/png, etc.

	if strings.HasSuffix(source, ".jpg") {
		fmt.Println("Converting the input jpg file to png file")
		os.Rename(source, destination)
		return nil
	}

	return fmt.Errorf("Unsupported source format: %s", source)
}

func updateConversionStatus(id int64, status string) {
	db.Model(&models.Conversion{ID: id}).Update("status", status)
}
