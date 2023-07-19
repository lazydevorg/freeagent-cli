package auth

import (
	"fmt"
	"net/http"
)

type CallbackHandler struct {
	state    string
	codeChan chan string
}

func NewCallbackHandler(state string) CallbackHandler {
	return CallbackHandler{state: state, codeChan: make(chan string)}
}

func (c *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer close(c.codeChan)
	defer w.(http.Flusher).Flush()

	callbackState := r.FormValue("state")
	if callbackState != c.state {
		http.Error(w, "Authentication failed: 'state' value is incorrect", http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, "Go back to your terminal.")
	c.codeChan <- r.FormValue("code")
}
