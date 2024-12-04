package cdn

import (
	"io"
	"serverless/cdn/storage"
)

// Command interface provides methods
// for uploading Content to the CDN.
type Command interface {
	// Upload uploads a single file with the given ID and Content.
	Upload(id string, file io.Reader) error
}

type CommandCDN struct {
	Storage *storage.StorageCDN
	Handler *CDNHandler
}

func (command *CommandCDN) Upload(id string, file io.Reader) error {
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = command.Handler.Upload(id, content)
	if err != nil {
		return err
	}
	return command.Storage.StoreFile(id, content)
}
