package api

import (
	"net/http"
)

type LoggingRoundTripper struct {
	Base http.RoundTripper
}

// Prepare custom client that using a logging transport.
func NewLoggingHTTPClient() *http.Client {
	client := *http.DefaultClient
	client.Transport = LoggingRoundTripper{http.DefaultTransport}

	return &client
}

func (r LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = r.Base.RoundTrip(req)

	return
}
