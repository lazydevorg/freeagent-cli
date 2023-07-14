package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	authURL     = "https://api.sandbox.freeagent.com/v2/approve_app"
	tokenURL    = "https://api.sandbox.freeagent.com/v2/token_endpoint"
	redirectURL = "http://localhost:8080/callback"
)

func Authenticate() *oauth2.Token {
	oAuthConfig := oAuthConfig()
	state := randomState()

	tokenChan := make(chan *oauth2.Token)
	callbackHandler := http.NewServeMux()
	callbackHandler.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		state2 := r.FormValue("state")
		if state2 != state {
			log.Fatalln("Authentication failed: 'randomState' value is incorrect")
		}

		code := r.FormValue("code")
		token, err := oAuthConfig.Exchange(r.Context(), code)
		if err != nil {
			log.Fatalln("Authentication failed:", err)
		}
		tokenChan <- token
	})

	url := oAuthConfig.AuthCodeURL(state)
	fmt.Printf("Click on the following URL and proceed with the login: %s\n", url)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: callbackHandler,
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalln("Error starting callback server:", err)
		}
	}()

	token := <-tokenChan
	fmt.Printf("Access Token: %+v\n", token)
	_ = server.Close()
	log.Println("Authentication completed")
	return token
}

func randomState() string {
	_, err := io.ReadFull(rand.Reader, make([]byte, 16))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString([]byte("randomState"))
}

func oAuthConfig() *oauth2.Config {
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
