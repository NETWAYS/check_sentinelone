package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

const DefaultTimeout = 5 * time.Second

type Client struct {
	AuthConfig    *clientcredentials.Config
	HTTPClient    *http.Client
	ManagementURL string
	AuthToken     string
}

func NewClient(url, token string, timeout time.Duration) (c *Client) {
	c = &Client{
		ManagementURL: url,
		AuthToken:     token,
	}

	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	c.HTTPClient = NewLoggingHTTPClient()
	c.HTTPClient.Timeout = timeout

	return
}

func (c *Client) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	// nolint: noctx
	// TODO Add context
	req, err = http.NewRequest(method, c.ManagementURL+"/web/api/"+url, body)
	if err != nil {
		err = fmt.Errorf("could not create http request: %w", err)
	}

	req.Header.Set("Authorization", "ApiToken "+c.AuthToken)

	return
}

func (c *Client) Do(req *http.Request) (res *http.Response, err error) {
	res, err = c.HTTPClient.Do(req) //nolint: gosec
	if err != nil {
		err = fmt.Errorf("HTTP request failed: %w", err)
	}

	return
}

type LoggingRoundTripper struct {
	Base http.RoundTripper
}

// NewLoggingHTTPClient prepares a custom client that using a logging transport.
func NewLoggingHTTPClient() *http.Client {
	client := *http.DefaultClient
	client.Transport = LoggingRoundTripper{http.DefaultTransport}

	return &client
}

func (r LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = r.Base.RoundTrip(req)

	return
}
