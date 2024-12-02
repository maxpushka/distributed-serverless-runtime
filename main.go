package main

import (
	"serverless/config"
	"serverless/router"
)

func main() {
	conf := config.New()
	router.Start(conf)
}
