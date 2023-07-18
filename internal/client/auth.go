package client

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/cache"
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

type CallbackHandler struct {
	oAuthConfig *oauth2.Config
	state       string
	tokenChan   chan *oauth2.Token
}

func NewCallbackHandler(oAuthConfig *oauth2.Config, state string) CallbackHandler {
	tokenChan := make(chan *oauth2.Token)
	return CallbackHandler{oAuthConfig: oAuthConfig, state: state, tokenChan: tokenChan}
}

func (c CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer close(c.tokenChan)
	callbackState := r.FormValue("state")
	if callbackState != c.state {
		http.Error(w, "Authentication failed: 'state' value is incorrect", http.StatusInternalServerError)
		return
	}

	code := r.FormValue("code")
	token, err := c.oAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, "Authentication successful. Go back to your terminal.")
	c.tokenChan <- token
}

func Authenticate() *oauth2.Token {
	oAuthConfig := oAuthConfig()
	state := randomState()

	callbackHandler := NewCallbackHandler(oAuthConfig, state)
	tokenChan := callbackHandler.tokenChan
	callbackServer := http.NewServeMux()
	callbackServer.Handle("/callback", callbackHandler)

	url := oAuthConfig.AuthCodeURL(state)
	fmt.Printf("Click on the following URL and proceed with the login: %s\n", url)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: callbackServer,
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalln("Error starting callback server:", err)
		}
	}()

	token := <-tokenChan
	if token == nil {
		log.Fatalln("Authentication failed")
	}

	err := storeToken(token)
	if err != nil {
		log.Fatalln("Error storing authentication data")
	}

	log.Println("Authentication completed")
	_ = server.Close()
	return token
}

func storeToken(token *oauth2.Token) error {
	return cache.SaveJson("auth", token)
}

func loadToken() (*oauth2.Token, error) {
	token := &oauth2.Token{}
	err := cache.LoadJson("auth", token)
	if err != nil {
		return nil, err
	}
	return token, nil
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
