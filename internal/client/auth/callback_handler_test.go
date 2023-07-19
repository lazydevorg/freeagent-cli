package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestCallbackHandlerSuccess(t *testing.T) {
	values := url.Values{"code": {"123"}, "state": {"123"}}.Encode()
	req, err := http.NewRequest("GET", "/callback?"+values, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := NewCallbackHandler("123")

	go handler.ServeHTTP(rr, req)

	code := <-handler.codeChan
	if code != "123" {
		t.Errorf("Unexpected auth code: '%s'", code)
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Unexpected status code: %d", rr.Code)
	}
	if rr.Body.String() != "Go back to your terminal." {
		t.Errorf("Unexpected body: '%s'", rr.Body.String())
	}
}

func TestCallbackHandlerWrongState(t *testing.T) {
	values := url.Values{"code": {"123"}, "state": {"321"}}.Encode()
	req, err := http.NewRequest("GET", "/callback?"+values, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := NewCallbackHandler("123")

	go handler.ServeHTTP(rr, req)

	select {
	case code := <-handler.codeChan:
		if code != "" {
			t.Errorf("Unexpected auth code: '%s'", code)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Code channel timed out")
	}

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status code: %d", rr.Code)
	}
	if rr.Body.String() != "Authentication failed: 'state' value is incorrect\n" {
		t.Errorf("Unexpected body: '%s'", rr.Body.String())
	}
}
