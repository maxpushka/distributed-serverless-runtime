package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// ComputeChecksum computes a combined checksum for multiple files.
func ComputeChecksum(file io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
