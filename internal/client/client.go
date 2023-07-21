package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/client/auth"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	u "net/url"
)

var baseUrl = "https://api.sandbox.freeagent.com/v2/"

type Client struct {
	Http        *http.Client
	tokenSource oauth2.TokenSource
}

var client *Client

func ClientSingleton() *Client {
	if client == nil {
		client = newClient(context.Background())
	}
	return client
}

func SaveToken() {
	if client != nil {
		token, err := client.tokenSource.Token()
		if err != nil {
			fmt.Println("Error storing token")
		}
		err = auth.StoreToken(token)
		if err != nil {
			fmt.Println("Error storing token")
		}
	}
}

func newClient(ctx context.Context) *Client {
	token, err := auth.LoadToken()
	if err != nil {
		log.Fatalln("Error loading authentication data")
	}
	tokenSource := auth.OAuthConfig().TokenSource(ctx, token)
	return &Client{
		Http:        oauth2.NewClient(ctx, tokenSource),
		tokenSource: tokenSource,
	}
}

func getRequest(url string) ([]byte, error) {
	c := ClientSingleton()
	url, err := u.JoinPath(baseUrl, url)
	if err != nil {
		return nil, err
	}
	response, err := c.Http.Get(url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func GetEntity[T any](url string, entityName string) (*T, error) {
	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	entityResponse := map[string]T{}
	err = json.Unmarshal(body, &entityResponse)
	if err != nil {
		return nil, err
	}
	var entity = entityResponse[entityName]
	return &entity, nil
}

func GetArray[T any](url string, groupName string) ([]T, error) {
	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	entityResponse := map[string][]T{}
	err = json.Unmarshal(body, &entityResponse)
	if err != nil {
		return nil, err
	}
	var entity = entityResponse[groupName]
	return entity, nil
}
