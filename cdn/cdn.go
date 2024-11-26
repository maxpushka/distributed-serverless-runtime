package cdn

import "io"

type CDN interface {
	Command
	Query
}

type Command interface {
	Upload(id string, content io.Reader) error
}

type Query interface {
	ReadFile(id string) (content io.Reader, checksum string, err error)
	Checksum(id string) (string, error)
}
