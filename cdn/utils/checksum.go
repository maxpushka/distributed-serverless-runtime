package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// ComputeChecksum computes a combined checksum for multiple files.
func ComputeChecksum(file io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
