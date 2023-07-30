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
	"strings"
)

var baseUrl = "https://api.freeagent.com/v2/"

type Client struct {
	Http        *http.Client
	tokenSource oauth2.TokenSource
}

var client *Client

func Singleton() *Client {
	if client == nil {
		client = newClient(context.Background())
	}
	return client
}

func SaveToken() {
	if client != nil {
		token, err := client.tokenSource.Token()
		if err != nil {
			fmt.Println("Error retrieving token: %w", err)
		}
		err = auth.StoreToken(token)
		if err != nil {
			fmt.Println("Error storing token: %w", err)
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

func getRequest(apiUrl string, params map[string]string) ([]byte, error) {
	url, err := getUrl(apiUrl, params)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}
	c := Singleton()
	response, err := c.Http.Do(request)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

const perPageDefault = 5

func getUrl(url string, params map[string]string) (*string, error) {
	if !strings.HasPrefix(url, "http") {
		url = baseUrl + url
	}
	parsedUrl, err := u.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %w", err)
	}
	query := parsedUrl.Query()
	query.Set("per_page", fmt.Sprintf("%d", perPageDefault))
	for key, value := range params {
		query.Set(key, value)
	}
	parsedUrl.RawQuery = query.Encode()
	stringUrl := parsedUrl.String()
	return &stringUrl, nil
}

func GetEntity[T any](url string, entityName string) (*T, error) {
	body, err := getRequest(url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	var entityResponse map[string]T
	err = json.Unmarshal(body, &entityResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	var entity = entityResponse[entityName]
	return &entity, nil
}

func GetCollection[T any](url string, groupName string, params map[string]string) ([]T, error) {
	body, err := getRequest(url, params)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	var entityResponse map[string][]T
	err = json.Unmarshal(body, &entityResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	var entity = entityResponse[groupName]
	return entity, nil
}
