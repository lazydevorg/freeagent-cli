package main

import (
	"context"
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/client"
)

func main() {
	//client.Authenticate()
	ctx := context.Background()
	c := client.NewClient(ctx)
	usersClient := c.NewUsersClient()
	defer c.Close()
	profile, err := usersClient.PersonalProfile()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", profile)
}
