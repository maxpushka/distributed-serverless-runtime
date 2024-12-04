package storage

import (
	"errors"
	"os"
	"path/filepath"
)

// Define the base directory for the CDN files.
const baseDir = "./distributed-cdn/"

// Ensure the base directory exists.
func init() {
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		panic("Failed to create base directory: " + err.Error())
	}
}

// StorageCDN provides methods for file operations in the CDN.
type StorageCDN struct{}

// StoreFile stores a file with the given ID in the CDN directory.
func (cdn *StorageCDN) StoreFile(id string, file []byte) error {
	filePath := filepath.Join(baseDir, id)
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = outFile.Write(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return nil
}

// RetrieveFile retrieves a file associated with the given ID.
func (cdn *StorageCDN) RetrieveFile(id string) (*os.File, error) {
	filePath := filepath.Join(baseDir, id)
	inFile, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return inFile, nil
}

// DeleteFile deletes a file associated with the given ID.
func (cdn *StorageCDN) DeleteFile(id string) error {
	filePath := filepath.Join(baseDir, id)
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not found")
		}
		return err
	}
	return nil
}
