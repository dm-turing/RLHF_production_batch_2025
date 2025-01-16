package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	fileID := vars["id"]

	// Replace with proper logic to retrieve file location
	dst, _ := strconv.Atoi(fileID)
	filePath := filepath.Join("uploads", uploadedfiles[dst])

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	http.ServeFile(w, r, filePath)
}
