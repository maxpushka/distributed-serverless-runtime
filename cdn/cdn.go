package cdn

// CDN defines the complete CDN interface
// with Command and Query capabilities.
type CDN interface {
	Command
	Query
}
