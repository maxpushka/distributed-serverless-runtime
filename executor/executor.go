package executor

import (
	"sync"
)

type Executor interface {
	Execute(id string, cfg map[string]string, req Request) (Response, error)
}

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    any               `json:"body,omitempty"`
}

type Response struct {
	StatusCode int               `json:"status_code,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

type PythonExecutor struct {
	// https://github.com/kluctl/go-embed-python
}

type GoExecutor struct {
	// https://github.com/traefik/yaegi
}

type Runner[Runtime any] struct {
	Mu sync.RWMutex
	// Checksum of the script source code.
	// Used to determine whether
	// the loaded script is up to date.
	Checksum string
	// Runtime for the script.
	Runtime Runtime
}
