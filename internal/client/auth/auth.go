package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/cache"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io"
	"log"
)

const (
	defaultAuthUrl     = "https://api.freeagent.com/v2/approve_app"
	defaultTokenUrl    = "https://api.freeagent.com/v2/token_endpoint"
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

type OAuthServer struct {
	config *oauth2.Config
	server *CallbackServer
}

func Authenticate(ctx context.Context, force bool) {
	if force {
		_ = DeleteToken()
	}
	NewOAuthServer().Authenticate(ctx)
}

func NewOAuthServer() *OAuthServer {
	config := OAuthConfig()
	return &OAuthServer{config: config, server: NewCallbackServer(config)}
}

func (s *OAuthServer) AuthCodeURL() string {
	return s.server.AuthCodeURL()
}

func (s *OAuthServer) Authenticate(ctx context.Context) *oauth2.Token {
	token, _ := LoadToken()
	if token != nil {
		return token
	}

	url := s.server.AuthCodeURL()
	fmt.Printf("Click on the following URL and proceed with the login: %s\n", url)

	code, err := s.server.WaitForAuthCode(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	token, err = s.config.Exchange(ctx, code)
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

func DeleteToken() error {
	return cache.Delete("auth")
}

func randomState() string {
	data := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, data)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(data)
}

func clientID() string {
	clientID := viper.GetString("auth.clientId")
	if clientID == "" {
		log.Fatal("CLIENT_ID environment variable must be set")
	}
	return viper.GetString("auth.clientId")
}

func clientSecret() string {
	clientSecret := viper.GetString("auth.clientSecret")
	if clientSecret == "" {
		log.Fatal("CLIENT_SECRET environment variable must be set")
	}
	return clientSecret
}

func authURL() string {
	return viper.GetString("auth.url")
}

func tokenURL() string {
	return viper.GetString("auth.tokenUrl")
}

func init() {
	_ = viper.BindEnv("auth.clientId", "CLIENT_ID")
	_ = viper.BindEnv("auth.clientSecret", "CLIENT_SECRET")
	_ = viper.BindEnv("auth.url", "AUTH_URL")
	viper.SetDefault("auth.url", defaultAuthUrl)
	_ = viper.BindEnv("auth.tokenUrl", "TOKEN_URL")
	viper.SetDefault("auth.tokenUrl", defaultTokenUrl)
}
