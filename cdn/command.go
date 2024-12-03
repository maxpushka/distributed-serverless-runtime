package cdn

import "io"

// Command interface provides methods
// for uploading content to the CDN.
type Command interface {
	// Upload uploads a single file with the given ID and content.
	Upload(id string, files []io.Reader) error
}
