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
	err = router.Start(conf)
	if err != nil {
		panic(err)
	}
}
