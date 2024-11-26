package js

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"serverless/cdn"
	"serverless/executor"

	"github.com/dop251/goja"
	"github.com/jellydator/ttlcache/v3"
)

type Executor struct {
	source  cdn.Query
	runners *ttlcache.Cache[string, *executor.Runner[*goja.Runtime]]
}

func NewExecutor(hot time.Duration, source cdn.Query) *Executor {
	runtimes := ttlcache.New(ttlcache.WithTTL[string, *executor.Runner[*goja.Runtime]](hot))
	go runtimes.Start() // starts automatic expired item deletion

	return &Executor{
		source:  source,
		runners: runtimes,
	}
}

func (exec *Executor) Execute(id string, cfg map[string]string, req executor.Request) (executor.Response, error) {
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
		runner.Mu.RUnlock() // unlock before acquiring write lock
		if err := exec.setupRunner(runner, sum, id); err != nil {
			return executor.Response{}, err
		}
		runner.Mu.RLock() // relock for reading
	}

	// Marshal request
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return executor.Response{}, err
	}
	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return executor.Response{}, err
	}

	// Execute the script
	outputObject, err := runner.Runtime.RunString(fmt.Sprintf("handle(%s, %s)", reqJSON, cfgJSON))
	if err != nil {
		return executor.Response{}, err
	}

	// Convert output to executor.Response
	var response executor.Response
	outputMap := outputObject.Export()
	outputBytes, err := json.Marshal(outputMap)
	if err != nil {
		return executor.Response{}, err
	}
	if err := json.Unmarshal(outputBytes, &response); err != nil {
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