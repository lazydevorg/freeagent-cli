package auth

import (
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

type CallbackServer struct {
	oAuthConfig *oauth2.Config
	state       string
}

func NewCallbackServer(oAuthConfig *oauth2.Config) *CallbackServer {
	return &CallbackServer{oAuthConfig: oAuthConfig, state: randomState()}
}

func (s *CallbackServer) AuthCodeURL() string {
	return s.oAuthConfig.AuthCodeURL(s.state)
}

func (s *CallbackServer) WaitForToken() *oauth2.Token {
	callbackHandler := NewCallbackHandler(s.oAuthConfig, s.state)
	callbackServer := http.NewServeMux()
	callbackServer.Handle("/callback", &callbackHandler)

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

	token := <-callbackHandler.tokenChan
	if token == nil {
		log.Fatalln("Authentication failed")
		return nil
	}

	log.Println("Authentication completed")
	_ = server.Close()
	return token
}
