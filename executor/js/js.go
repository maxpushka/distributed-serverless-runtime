package js

import (
	"encoding/json"
	"fmt"
	"io"

	"serverless/cdn"
	"serverless/executor"

	"github.com/dop251/goja"
	"github.com/jellydator/ttlcache/v3"
)

type Executor struct {
	source  cdn.Query
	runners *ttlcache.Cache[string, *executor.Runner[*goja.Runtime]]
}

func NewExecutor(source cdn.Query) *Executor {
	runtimes := ttlcache.New(ttlcache.WithTTL[string, *executor.Runner[*goja.Runtime]](executor.HotDuration))
	go runtimes.Start() // starts automatic expired item deletion

	return &Executor{
		source:  source,
		runners: runtimes,
	}
}

func (exec *Executor) Execute(id string, req executor.Request) (executor.Response, error) {
	// TODO:
	// 1. Get or create runtime, load the script.
	//    For a given user runtime may be:
	//      - hot in cache with the script already loaded
	//      - cold, need to fetch the script
	// 2. Execute the script
	// 3. Return the formatted response

	sum, err := exec.source.Checksum(id)
	if err != nil {
		return executor.Response{}, err
	}

	hotRunner, ok := exec.runners.GetOrSet(id, &executor.Runner[*goja.Runtime]{})
	runner := hotRunner.Value()

	// Check whether script and runtime are up to date
	runner.Mu.RLock()
	defer runner.Mu.RUnlock()
	if !ok || sum != runner.Checksum { // script was updated, need to reload
		if err := exec.setupRunner(runner, sum, id); err != nil {
			return executor.Response{}, err
		}
	}

	// Execute the script
	request, err := json.Marshal(req)
	if err != nil {
		return executor.Response{}, err
	}
	output, err := runner.Runtime.RunString(fmt.Sprintf("handle(%s)", request))
	if err != nil {
		return executor.Response{}, err
	}

	// Convert output to executor.Response
	var response executor.Response
	outputStr := output.String()
	if err := json.Unmarshal([]byte(outputStr), &response); err != nil {
		return executor.Response{}, fmt.Errorf("failed to unmarshal handler response: %w", err)
	}

	return response, nil
}

func (exec *Executor) setupRunner(runner *executor.Runner[*goja.Runtime], prevSum, id string) error {
	runner.Mu.Lock()
	defer runner.Mu.Unlock()

	// Verify again the checksum differs
	// in case another goroutine updated the script
	// while we were waiting for the lock
	if prevSum == runner.Checksum {
		return nil
	}

	reader, sum, err := exec.source.ReadFile(id)
	if err != nil {
		return err
	}
	source, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	if _, err := vm.RunString(string(source)); err != nil {
		return err
	}

	runner.Checksum = sum
	runner.Runtime = vm
	return nil
}
