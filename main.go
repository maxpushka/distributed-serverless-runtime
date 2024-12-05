package main

import (
	"context"
	"serverless/cdn"
	"serverless/cdn/storage"
	"serverless/config"
	"serverless/executor/js"
	"serverless/router"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	storageCDN := &storage.StorageCDN{}

	handler, err := cdn.InitCDNHandler(ctx, storageCDN, true)
	if err != nil {
		panic(err)
	}

	command := &cdn.CommandCDN{
		Storage: storageCDN,
		Handler: handler,
	}
	query := cdn.QueryCDN{
		Storage: storageCDN,
	}
	runtime := js.NewExecutor(conf.Executor, query)
	router.Start(conf, command, runtime)
}
