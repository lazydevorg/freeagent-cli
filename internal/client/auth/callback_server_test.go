package auth

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

func TestTokenExchange(t *testing.T) {
	config := oauth2.Config{}

	tests := []struct {
		name           string
		serverState    string
		requestState   string
		requestCode    string
		expectedStatus int
	}{
		{"Success", "123", "123", "321", http.StatusOK},
		{"Wrong state", "123", "321", "", http.StatusInternalServerError},
		{"No code", "123", "123", "", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				server := NewCallbackServer(&config)
				server.state = tt.serverState
				code, _ := server.WaitForAuthCode()
				if code != tt.requestCode {
					t.Error("Unexpected auth code:", code)
				}
			}()

			values := url.Values{"code": {tt.requestCode}, "state": {tt.requestState}}.Encode()
			req, err := http.NewRequest("GET", "http://localhost:8080/callback?"+values, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			var client http.Client
			response, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != tt.expectedStatus {
				t.Error("Unexpected status code:", response.StatusCode)
			}

			wg.Wait()
			fmt.Println("Done")
		})
	}
}

func TestAuthURL(t *testing.T) {
	config := oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://oauth.test/auth",
		},
		Scopes: []string{"scope"},
	}
	server := NewCallbackServer(&config)
	if server.AuthCodeURL() != "https://oauth.test/auth?client_id=CLIENT_ID&response_type=code&scope=scope&state=cmFuZG9tU3RhdGU%3D" {
		t.Error("Unexpected auth URL:", server.AuthCodeURL())
	}
}
