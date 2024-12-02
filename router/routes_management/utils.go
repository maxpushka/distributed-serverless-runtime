package routes_management

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func ParseRouteID(r *http.Request) (string, int, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return "", 0, err
	}
	return idStr, id, nil
}

func SaveFile(w http.ResponseWriter, r *http.Request, idStr string, baseDir string) error {
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 MB max
	err := r.ParseMultipartForm(10 << 20)           // 10 MB max
	if err != nil {
		log.Print(err)
		return err
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Print(err)
		return err
	}
	defer file.Close()

	filename := filepath.Base(handler.Filename)
	filename = idStr + filepath.Ext(filename)

	dst, err := os.Create(filepath.Join(baseDir, filename))
	if err != nil {
		log.Print(err)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func FindFileByName(idStr string, baseDir string) (string, error) {
	var foundFile string
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			baseName := filepath.Base(path)
			fileNameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
			if fileNameWithoutExt == idStr {
				foundFile = path
				return errors.New("file found") // Use error to break out of WalkDir
			}
		}
		return nil
	})
	if err != nil && err.Error() != "file found" {
		return "", err
	}
	if foundFile == "" {
		return "", errors.New("file not found")
	}
	return foundFile, nil
}

func ServeFile(w http.ResponseWriter, filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the first 512 bytes to detect the content type
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Reset the read pointer to the beginning of the file
	file.Seek(0, io.SeekStart)

	// Detect the content type
	contentType := http.DetectContentType(buf[:n])

	// Set the headers
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))

	// Write the file content to the response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving file", http.StatusInternalServerError)
		return
	}
}
