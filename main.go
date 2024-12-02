package main

import (
	"serverless/config"
	"serverless/router"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	router.Start(conf)
}
