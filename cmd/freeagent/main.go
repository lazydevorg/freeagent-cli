package main

import (
	"github.com/lazydevorg/freeagent-cli/internal/client/auth"
)

func main() {
	auth.NewOAuthServer().Authenticate()
	//ctx := context.Background()
	//c := client.NewClient(ctx)
	//usersClient := c.NewUsersClient()
	//defer c.Close()
	//profile, err := usersClient.PersonalProfile()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v", profile)
}
