package client

import (
	"context"
	"encoding/json"
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

func (c *Client) Close() error {
	token, err := c.tokenSource.Token()
	if err != nil {
		return err
	}
	return storeToken(token)
}

func NewClient(ctx context.Context) *Client {
	token, err := loadToken()
	if err != nil {
		log.Fatalln("Error loading authentication data")
	}
	tokenSource := oAuthConfig().TokenSource(ctx, token)
	return &Client{
		Http:        oauth2.NewClient(ctx, tokenSource),
		tokenSource: tokenSource,
	}
}

func GetEntity[T any](c *Client, url string, entityName string) (*T, error) {
	url, err := u.JoinPath(baseUrl, url)
	if err != nil {
		return nil, err
	}
	response, err := c.Http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
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

func GetArray[T any](c *Client, url string, groupName string) ([]T, error) {
	url, err := u.JoinPath(baseUrl, url)
	if err != nil {
		return nil, err
	}
	response, err := c.Http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
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
