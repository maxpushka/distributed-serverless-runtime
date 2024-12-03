package cdn

import "io"

// CDN defines the complete CDN interface
// with Command and Query capabilities.
type CDN interface {
	Command
	Query
}

// Command interface provides methods
// for uploading content to the CDN.
type Command interface {
	// Upload uploads a single file with the given ID and content.
	Upload(id string, files []io.Reader) error
}

// Query interface provides methods
// for reading and checking files in the CDN.
type Query interface {
	// ReadFiles retrieves multiple files associated with a given ID.
	ReadFiles(id string) (files []io.Reader, checksum string, err error)

	// Checksum retrieves a single checksum
	// for all files associated with the given ID.
	// It computes individual checksums,
	// sorts them, concatenates them, and hashes the result.
	Checksum(id string) (string, error)
}
