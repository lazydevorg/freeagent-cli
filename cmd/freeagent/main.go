package main

import (
	"context"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"io"
)

func main() {
	//client.Authenticate()
	ctx := context.Background()
	client := client.NewClient(ctx)
	defer client.Close()
	response, err := client.Http.Get("https://api.sandbox.freeagent.com/v2/company")
	if err != nil {
		return
	}
	defer response.Body.Close()
	_, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
}
