package client

import (
	"context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

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
