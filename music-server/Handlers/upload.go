package handlers

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Function to handle Validation of the uploaded file so that only music files are sent
func ValidateFileType(file []byte) (string, error) {
	filetype := http.DetectContentType(file)
	switch filetype {
	case "audio/mpeg", "audio/ogg", "audio/flac", "audio/wav", "application/zip":
		return filetype, nil
	default:
		return "", fmt.Errorf("unsupported file type: %v", filetype)
	}
}

// The function to handle extraction of zipped folders
func ExtractZip(string) (string, error) {
	// Unzip the file
	// Create a directory to extract to
	// Use the archive/zip package to unzip
	return "", nil

}

func TrackUploader(w http.ResponseWriter, r *http.Request) {

	// Set a maximum upload size of 50MB
	const maxUpload = 100 << 20

	// The form field name for the file upload
	const fileField = "myfile"
	//const temp = handler.filename

	// Limit the size of the request body to prevent large file uploads
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)

	// Parse the multipart form in the request
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}
	// Retrieve the file from form data
	file, handler, err := r.FormFile(fileField)
	if err != nil {
		http.Error(w, "Error Retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//fmt.Fprintf(w, "Uploaded File: %+v\n", handler.Filename)
	//fmt.Fprintf(w, "File size: %d bytes\n", handler.Size)
	//fmt.Fprintf(w, "Uploaded Mime: %v\n", handler.Header)
	// log the file details
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File size: %d bytes\n", handler.Size)
	log.Printf("Uploaded Mime: %v\n", handler.Header)

	// Read the file content into a byte slice
	if f, err := io.ReadAll(file); err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	} else {
		// Validate the file type
		if val, err := ValidateFileType(f); err != nil {
			http.Error(w, fmt.Sprintf("Invalid file type: %v", err), http.StatusBadRequest)
			return
		} else {
			fmt.Printf("File type is valid: %v\n", val)
		}
	}
	// Reset the file pointer to the beginning of the file for further operations
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, "Error resetting file pointer", http.StatusInternalServerError)
		return
	}

	// Create a directory if not existing
	dir_path := os.Getenv("LOCAL_STORAGE_PATH")
	if _, err := os.Stat(dir_path); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %v\n", err)
		if err := os.MkdirAll(dir_path, 0775); err != nil {
			fmt.Printf("Unable to create directory: %v\n", err)
			return
		}
	}
	// Check if the uploaded file is a zip file or an audio file based on its MIME type
	switch handler.Header.Get("Content-Type") {
	case "application/zip":
		// Handle zip file extraction
		log.Println("Zip file detected, proceeding to extract...")
		// save the uploaded zip file to a temporary location
		tempZipPath := filepath.Join(dir_path, "temp_"+handler.Filename)

		// Create the temporary file

		var tmpzip *os.File

		if tmpzip, err = os.Create(tempZipPath); err != nil {
			http.Error(w, "Unable to create temp zip file", http.StatusInternalServerError)
			return
		}
		defer tmpzip.Close()

		//defer os.Remove(tempZipPath) // Clean up the temp file after extraction

		// Copy the uploaded zip file to the temporary file
		if _, err := io.Copy(tmpzip, file); err != nil {
			http.Error(w, "Unable to save temp zip file", http.StatusInternalServerError)
			return
		}
		tmpzip.Close() // Close the file to ensure all data is written(important for Windows)

		// Open the zip file
		archive1, err := zip.OpenReader(tempZipPath)
		if err != nil {
			http.Error(w, "Error opening zip file", http.StatusInternalServerError)
			return
		}
		defer archive1.Close() // Close the zip file when done

		// Iterate through the files in the zip archive
		for _, f := range archive1.File {
			dstForALbum := filepath.Join(dir_path, "albums", handler.Filename, f.Name)
			fmt.Println("unzipping file ", dstForALbum)

			// Ensure the destination path is within the intended directory to prevent Zip Slip vulnerability
			if !strings.HasPrefix(dstForALbum, filepath.Clean(dir_path)+string(os.PathSeparator)) {
				fmt.Println("invalid file path")
				return
			}
			// Create directories as needed and extract files
			if f.FileInfo().IsDir() {
				fmt.Println("creating directory...")
				os.MkdirAll(dstForALbum, os.ModePerm)
				continue
			}
			// Create the necessary directories for the file
			if err := os.MkdirAll(dstForALbum, os.ModePerm); err != nil {
				panic(err)
			}
			// Create the destination file
			dstForAlbumAg, err := os.OpenFile(dstForALbum, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				panic(err)
			}
			// Open the file inside the zip archive
			archive, err := f.Open()
			if err != nil {
				http.Error(w, "Error opening file in zip", http.StatusInternalServerError)
				return
			}
			defer archive.Close() // Close the file when done
			// Copy the file content to the destination file
			if _, err := io.Copy(dstForAlbumAg, archive); err != nil {
				panic(err)
			}
			defer dstForAlbumAg.Close()
			// Respond to the client that the upload was successful

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Folder uploaded successfully!"))
			//json.NewEncoder(w).Encode(handler.Filename)
		}
		return
	case "audio/mpeg", "audio/ogg", "audio/flac", "audio/wav":
		// Proceed to save the audio file
		log.Println("Audio file detected, proceeding to save...")
		// Create the full path to the file
		dstPath := filepath.Join(dir_path, handler.Filename)
		// Create destination file making sure the path is writeable.
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Unable to create the file ", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		// Copy the uploaded file to the destination file
		if _, err = io.Copy(dst, file); err != nil {
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
			return
		}
		// Respond to the client that the upload was successful
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully!"))
		json.NewEncoder(w).Encode(handler.Filename)
	default: // Unsupported file type
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("File uploaded successfully!"))
	//json.NewEncoder(w).Encode(handler.Filename)

}
