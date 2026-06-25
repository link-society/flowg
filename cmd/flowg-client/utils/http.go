package utils

import (
	"fmt"

	"net/http"
	"net/url"
)

// Client is a minimal HTTP client that targets a FlowG API base URL and
// authenticates every request with a bearer token.
type Client struct {
	BaseUrl    string
	Token      string
	httpClient *http.Client
}

// NewClient builds a Client for the given API base URL and bearer token. An
// empty token leaves requests unauthenticated.
func NewClient(baseUrl string, token string) *Client {
	return &Client{
		BaseUrl:    baseUrl,
		Token:      token,
		httpClient: &http.Client{},
	}
}

// Do resolves the request path against the configured base URL, attaches the
// bearer token, and sends the request.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	fullUrl, err := url.JoinPath(c.BaseUrl, req.URL.Path)
	if err != nil {
		return nil, err
	}

	queryset := req.URL.Query()
	req.URL, err = url.Parse(fullUrl)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = queryset.Encode()

	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	return c.httpClient.Do(req)
}
