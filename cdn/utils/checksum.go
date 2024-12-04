package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeChecksum computes a combined checksum for multiple files.
func ComputeChecksum(file []byte) (string, error) {
	hash := sha256.New()
	hash.Write(file)

	return hex.EncodeToString(hash.Sum(nil)), nil
}
