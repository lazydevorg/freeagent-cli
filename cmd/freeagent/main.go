package main

import (
	"context"
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/auth"
	"io"
)

func main() {
	//auth.Authenticate()
	client := auth.NewClient(context.Background())
	response, err := client.Get("https://api.sandbox.freeagent.com/v2/company")
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Println(string(data))
}
