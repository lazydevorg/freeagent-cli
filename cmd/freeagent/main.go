package main

import (
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/cache"
	"golang.org/x/oauth2"
)

func main() {
	//auth.Authenticate()
	var Token oauth2.Token
	cache.LoadJson("auth", &Token)
	fmt.Println(Token)
}
