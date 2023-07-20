package auth

import (
	"context"
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
	defaultAuthUrl     = "https://api.sandbox.freeagent.com/v2/approve_app"
	defaultTokenUrl    = "https://api.sandbox.freeagent.com/v2/token_endpoint"
	defaultRedirectUrl = "http://localhost:8080/callback"
)

func OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID(),
		ClientSecret: clientSecret(),
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL(),
			TokenURL:  tokenURL(),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: defaultRedirectUrl,
	}
}

func Authenticate() *oauth2.Token {
	oAuthConfig := OAuthConfig()
	server := NewCallbackServer(oAuthConfig)

	url := server.AuthCodeURL()
	fmt.Printf("Click on the following URL and proceed with the login: %s\n", url)

	code, err := server.WaitForAuthCode()
	if err != nil {
		log.Fatalln(err)
	}

	token, err := oAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalln("Authentication failed:", err)
		return nil
	}

	err = StoreToken(token)
	if err != nil {
		log.Fatalln("Error storing authentication data")
	}
	fmt.Println("Authentication successful")

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
	data := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, data)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(data)
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

func authURL() string {
	value := os.Getenv("AUTH_URL")
	if value == "" {
		return defaultAuthUrl
	}
	return value
}

func tokenURL() string {
	value := os.Getenv("TOKEN_URL")
	if value == "" {
		return defaultTokenUrl
	}
	return value
}
