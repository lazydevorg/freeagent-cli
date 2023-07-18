package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/cache"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
)

const (
	authURL     = "https://api.sandbox.freeagent.com/v2/approve_app"
	tokenURL    = "https://api.sandbox.freeagent.com/v2/token_endpoint"
	redirectURL = "http://localhost:8080/callback"
)

func OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID(),
		ClientSecret: clientSecret(),
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: redirectURL,
	}
}

func Authenticate() *oauth2.Token {
	oAuthConfig := OAuthConfig()
	server := NewCallbackServer(oAuthConfig)

	url := server.AuthCodeURL()
	fmt.Printf("Click on the following URL and proceed with the login: %s\n", url)

	token := server.WaitForToken()
	if token == nil {
		log.Fatalln("Authentication failed")
	}

	err := StoreToken(token)
	if err != nil {
		log.Fatalln("Error storing authentication data")
	}

	return token
}

func LoadToken() (*oauth2.Token, error) {
	token := &oauth2.Token{}
	err := cache.LoadJson("auth", token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func StoreToken(token *oauth2.Token) error {
	return cache.SaveJson("auth", token)
}

func randomState() string {
	_, err := io.ReadFull(rand.Reader, make([]byte, 16))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString([]byte("randomState"))
}

func clientSecret() string {
	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("CLIENT_SECRET environment variable must be set")
	}
	return clientSecret
}

func clientID() string {
	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		log.Fatal("CLIENT_ID environment variable must be set")
	}
	return clientID
}
