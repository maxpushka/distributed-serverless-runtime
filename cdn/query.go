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

	// Checksum retrieves a single checksum
	// for the file associated with the given ID.
	// It computes individual checksum and hashes it
	Checksum(id string) (string, error)
}

type QueryCDN struct {
	storage storage.StorageCDN
}

func (cdn *QueryCDN) ReadFile(id string) (file io.Reader, checksum string, err error) {
	file, err = cdn.storage.RetrieveFile(id)
	if err != nil {
		return nil, "", err
	}

	checksum, err = utils.ComputeChecksum(file)
	if err != nil {
		return nil, "", err
	}

	return file, checksum, nil
}

func (cdn *QueryCDN) Checksum(id string) (string, error) {
	file, err := cdn.storage.RetrieveFile(id)
	if err != nil {
		return "", err
	}

	checksum, err := utils.ComputeChecksum(file)
	if err != nil {
		return "", err
	}

	return checksum, nil
}
