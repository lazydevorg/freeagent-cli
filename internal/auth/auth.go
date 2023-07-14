package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
			http.Error(w, "Authentication failed: 'state' value is incorrect", http.StatusInternalServerError)
			close(tokenChan)
			return
		}

		code := r.FormValue("code")
		token, err := oAuthConfig.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Authentication failed: "+err.Error(), http.StatusInternalServerError)
			close(tokenChan)
			return
		}
		_, _ = fmt.Fprintf(w, "Authentication successful. Go back to your terminal.")
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
	jsonToken, err := json.Marshal(token)
	if err != nil {
		log.Fatalln("Error storing authentication data:", err)
	}
	base64Token := base64.StdEncoding.EncodeToString(jsonToken)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cachePath := filepath.Join(home, ".cache", "freeagent-cli")
	err = os.MkdirAll(cachePath, 0700)
	if err != nil {
		return err
	}

	authDataPath := filepath.Join(cachePath, "auth")
	f, err := os.Create(authDataPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(base64Token)
	if err != nil {
		return err
	}
	return nil
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
