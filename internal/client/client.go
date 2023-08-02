package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/client/auth"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	u "net/url"
	"reflect"
	"strings"
)

var baseUrl = "https://api.freeagent.com/v2/"

type Client struct {
	Http        *http.Client
	BaseUrl     string
	tokenSource oauth2.TokenSource
}

var client *Client

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

func NewClient(ctx context.Context) *Client {
	token, err := auth.LoadToken()
	if err != nil {
		log.Fatalln("Error loading authentication data")
	}
	tokenSource := auth.OAuthConfig().TokenSource(ctx, token)
	return &Client{
		Http:        oauth2.NewClient(ctx, tokenSource),
		BaseUrl:     baseUrl,
		tokenSource: tokenSource,
	}
}

func (c *Client) getRequest(apiUrl string, params map[string]string) ([]byte, error) {
	url, err := c.getUrl(apiUrl, params)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.Http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

const perPageDefault = 5

func (c *Client) getUrl(url string, params map[string]string) (*string, error) {
	if !strings.HasPrefix(url, "http") {
		url = c.BaseUrl + url
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

func GetEntity[T any](c *Client, url string, entityName string) (*T, error) {
	body, err := c.getRequest(url, nil)
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

func GetCollection[T any](c *Client, url string, groupName string, params map[string]string) ([]T, error) {
	body, err := c.getRequest(url, params)
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

func (c *Client) PostRequest(apiUrl string, data []byte) (*http.Response, error) {
	url, err := c.getUrl(apiUrl, nil)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", *url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return c.Http.Do(request)
}

func PostEntity[T any](c *Client, url string, name string, entity *T) (*T, error) {
	data := map[string]*T{name: entity}
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling entity: %w", err)
	}
	response, err := c.PostRequest(url, encoded)
	fmt.Println(response)
	return nil, nil
}

type RelatedEntityProcessorFunc func(entity map[string]interface{}) string

func GetRelatedEntities[E any](c *Client, entities []E, entityFieldName string, related map[string]string, processorFunc RelatedEntityProcessorFunc) error {
	for _, entity := range entities {
		v := reflect.ValueOf(entity)
		field := v.FieldByName(entityFieldName)
		url := field.String()
		if _, found := related[url]; !found {
			re, err := GetEntity[map[string]interface{}](c, url, strings.ToLower(entityFieldName))
			if err != nil {
				return fmt.Errorf("error getting related entity: %w", err)
			}
			related[url] = processorFunc(*re)
		}
	}
	return nil
}
