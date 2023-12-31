package auth

import (
	"errors"
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

func (s *CallbackServer) WaitForAuthCode() (string, error) {
	callbackHandler := NewCallbackHandler(s.state)
	callbackServer := http.NewServeMux()
	callbackServer.Handle("/callback", &callbackHandler)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: callbackServer,
	}
	defer closeServer(&server)

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalln("Error starting callback server:", err)
		}
	}()

	code := <-callbackHandler.codeChan
	if code == "" {
		return "", errors.New("authentication failed")
	}

	return code, nil
}

func closeServer(server *http.Server) {
	err := server.Close()
	if err != nil {
		log.Fatalln("Error closing callback server:", err)
	}
}
