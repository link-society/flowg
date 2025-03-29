package client

import (
	"fmt"

	"net/http"
	"net/url"
)

type Client struct {
	BaseUrl    string
	Token      string
	httpClient *http.Client
}

func NewClient(baseUrl string, token string) *Client {
	return &Client{
		BaseUrl:    baseUrl,
		Token:      token,
		httpClient: &http.Client{},
	}
}

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
