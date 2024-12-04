package cdn

import (
	"io"
	"serverless/cdn/storage"
	"serverless/cdn/utils"
)

// Query interface provides methods
// for reading and checking files in the CDN.
type Query interface {
	// ReadFile retrieves  file associated with a given ID.
	ReadFile(id string) (file io.Reader, checksum string, err error)

	// Checksum retrieves a single Checksum
	// for the file associated with the given ID.
	// It computes individual Checksum and hashes it
	Checksum(id string) (string, error)
}

type QueryCDN struct {
	Storage *storage.StorageCDN
}

func (cdn *QueryCDN) ReadFile(id string) (content []byte, checksum string, err error) {
	file, err := cdn.Storage.RetrieveFile(id)
	if err != nil {
		return nil, "", err
	}

	content, err = io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	checksum, err = utils.ComputeChecksum(content)
	if err != nil {
		return nil, "", err
	}

	return content, checksum, nil
}

func (cdn *QueryCDN) Checksum(id string) (string, error) {
	file, err := cdn.Storage.RetrieveFile(id)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	checksum, err := utils.ComputeChecksum(content)
	if err != nil {
		return "", err
	}

	return checksum, nil
}
