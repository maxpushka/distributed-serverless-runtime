package executor

import (
	"sync"
	"time"
)

const HotDuration = 30 * time.Minute

type Executor interface {
	Execute(id string, req Request) (Response, error)
}

type Request struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    []byte              `json:"body,omitempty"`
}

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

type PythonExecutor struct {
	// https://github.com/kluctl/go-embed-python
}

type GoExecutor struct {
	// https://github.com/traefik/yaegi
}

type Runner[Runtime any] struct {
	Mu       sync.RWMutex
	Checksum string
	Runtime  Runtime
}
