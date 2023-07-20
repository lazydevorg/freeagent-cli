package auth

import (
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestOAuthConfig(t *testing.T) {
	t.Setenv("CLIENT_ID", "client_id")
	t.Setenv("CLIENT_SECRET", "client_secret")

	config := OAuthConfig()

	if config.ClientID != "client_id" {
		t.Error("Unexpected client ID:", config.ClientID)
	}
	if config.ClientSecret != "client_secret" {
		t.Error("Unexpected client secret:", config.ClientSecret)
	}
	if config.Endpoint.AuthURL != "https://api.sandbox.freeagent.com/v2/approve_app" {
		t.Error("Unexpected auth URL:", config.Endpoint.AuthURL)
	}
	if config.Endpoint.TokenURL != "https://api.sandbox.freeagent.com/v2/token_endpoint" {
		t.Error("Unexpected token URL:", config.Endpoint.TokenURL)
	}
	if config.RedirectURL != "http://localhost:8080/callback" {
		t.Error("Unexpected redirect URL:", config.RedirectURL)
	}
}

func TestTokenStorage(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	token := &oauth2.Token{
		AccessToken:  "accesstoken",
		TokenType:    "bearer",
		RefreshToken: "refreshtoken",
		Expiry:       time.Now(),
	}

	err := StoreToken(token)
	if err != nil {
		t.Fatal(err)
	}

	loadedToken, err := LoadToken()
	if err != nil {
		t.Fatal(err)
	}
	if !sameToken(loadedToken, token) {
		t.Errorf("Loaded token %+v differs from the stored one %+v", loadedToken, token)
	}
}

func sameToken(t1, t2 *oauth2.Token) bool {
	return t1.AccessToken == t2.AccessToken &&
		t1.TokenType == t2.TokenType &&
		t1.RefreshToken == t2.RefreshToken &&
		t1.Expiry.Equal(t2.Expiry)
}

func TestAuthenticate(t *testing.T) {
	var wg sync.WaitGroup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "client_id=client_id&client_secret=client_secret&code=CODE&grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback" {
			t.Error("Unexpected body:", string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "{\"access_token\":\"accesstoken\",\"token_type\":\"bearer\",\"expires_in\":3600,\"refresh_token\":\"refreshtoken\"\n\n}")
		w.(http.Flusher).Flush()
	}))
	defer ts.Close()

	t.Setenv("AUTH_URL", ts.URL+"/auth")
	t.Setenv("TOKEN_URL", ts.URL+"/token")
	t.Setenv("CLIENT_ID", "client_id")
	t.Setenv("CLIENT_SECRET", "client_secret")

	wg.Add(1)
	oAuthServer := NewOAuthServer()
	go func() {
		defer wg.Done()
		token := oAuthServer.Authenticate()
		if token == nil {
			t.Error("Unexpected nil token")
		}
	}()

	state := authCodeState(oAuthServer.AuthCodeURL())
	values := url.Values{"code": {"CODE"}, "state": {state}}.Encode()
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
	if response.StatusCode != http.StatusOK {
		t.Error("Unexpected status code:", response.StatusCode)
	}

	wg.Wait()
}

func authCodeState(authCodeURL string) string {
	u, _ := url.Parse(authCodeURL)
	vs, _ := url.ParseQuery(u.RawQuery)
	return vs.Get("state")
}
