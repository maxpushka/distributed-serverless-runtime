package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"serverless/cdn"
	"serverless/cdn/storage"
	"strings"
)

func main() {
	ctx := context.Background()
	storageCDN := &storage.StorageCDN{}

	handler, err := cdn.InitCDNHandler(ctx, storageCDN, true)
	if err != nil {
		panic(err)
	}

	query := &cdn.QueryCDN{
		Storage: storageCDN,
	}

	command := &cdn.CommandCDN{
		Storage: storageCDN,
		Handler: handler,
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		content, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if content == "exit" {
			break
		}

		params := strings.Split(strings.TrimSuffix(content, "\n"), " ")
		if len(params) < 2 {
			fmt.Println("Invalid command")
			continue
		}

		switch params[0] {
		case "upload":
			file := strings.NewReader(params[1])
			err = command.Upload(params[1], file)
		case "read":
			content, checksum, err := query.ReadFile(params[1])
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(checksum)
			fmt.Println(string(content))
		}
	}

	fmt.Println(query)
	fmt.Println(command)
}
