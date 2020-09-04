package api

import (
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"time"
)

const DefaultTimeout = 5

type Client struct {
	AuthConfig    *clientcredentials.Config
	HttpClient    *http.Client
	ManagementURL string
	AuthToken     string
}

func NewClient(url, token string) (c *Client) {
	c = &Client{
		ManagementURL: url,
		AuthToken:     token,
	}

	c.HttpClient = NewLoggingHTTPClient()
	c.HttpClient.Timeout = time.Duration(DefaultTimeout) * time.Second

	return
}

func (c *Client) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, c.ManagementURL+"/web/api/"+url, body)
	if err != nil {
		err = fmt.Errorf("could not create http request: %w", err)
	}

	req.Header.Set("Authorization", "ApiToken "+c.AuthToken)

	return
}

func (c *Client) Do(req *http.Request) (res *http.Response, err error) {
	res, err = c.HttpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("HTTP request failed: %w", err)
	}

	return
}
