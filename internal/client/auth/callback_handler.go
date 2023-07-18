package auth

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
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

func (c *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
