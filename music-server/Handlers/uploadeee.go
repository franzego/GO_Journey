/*
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)



func UploadHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//fmt.Println(r.Header)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Unable to Parse File", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["myfile"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	//Retreiving file from form input
	file1, handler, err := r.FormFile("myfile")
	if err != nil {
		http.Error(w, "Error Retrieving the file", http.StatusBadRequest)
		return
	}
	defer file1.Close()

	// Create a directory if not existing
	dir_path := os.Getenv("LOCAL_STORAGE_PATH")
	if _, err := os.Stat(dir_path); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %v\n", err)
		if err := os.MkdirAll(dir_path, 0775); err != nil {
			fmt.Printf("Unable to create directory: %v\n", err)
			return
		}
	}

	// We neeed to read the file to byte that can be used for the validation as the validation function accepts bytes
	f, err := os.ReadFile("myfile")
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	//Validation of sent file from client

	_, err = ValidateFileType(f)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid file type: %v", err), http.StatusBadRequest)
		return
	}

	// If the file is valid, we can proceed to save it or process it as needed
	fmt.Fprintf(w, "Uploaded File: %+v\n", handler.Filename)
	fmt.Fprintf(w, "File size: %d bytes\n", handler.Size)
	fmt.Fprintf(w, "Uploaded Mime: %v\n", handler.Header)

	// Save the uploaded file to the local storage path

	var Dst *os.File

	DstPath := filepath.Join(dir_path + "/" + handler.Filename)

	// Create the file on disk
	if Dst, err = os.Create(DstPath); err != nil {
		http.Error(w, "Unable to create the file for writing.", http.StatusInternalServerError)
		return
	}
	defer Dst.Close()

	//We need to copy the file to the new folder
	if _, err = io.Copy(Dst, file1); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully", "filename": handler.Filename})
}
*/package handlers
