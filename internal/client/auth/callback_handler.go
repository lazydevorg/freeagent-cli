package auth

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

type CallbackHandler struct {
	oAuthConfig *oauth2.Config
	state       string
	codeChan    chan string
}

func NewCallbackHandler(oAuthConfig *oauth2.Config, state string) CallbackHandler {
	return CallbackHandler{oAuthConfig: oAuthConfig, state: state, codeChan: make(chan string)}
}

func (c *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer close(c.codeChan)
	callbackState := r.FormValue("state")
	if callbackState != c.state {
		http.Error(w, "Authentication failed: 'state' value is incorrect", http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, "Go back to your terminal.")
	w.(http.Flusher).Flush()
	c.codeChan <- r.FormValue("code")
}
