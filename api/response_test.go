package api_test

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"strings"
	"testing"
)

const testData = `{
    "data": {
		"id": "c86dc437-daa8-4dec-b5a8-ce1e0a2c0c5e",
		"name": "test1"
    }
}`

const testResults = `{
	"data" : [
		{"id": "c86dc437-daa8-4dec-b5a8-ce1e0a2c0c5e", "name": "test1"}
	],
    "pagination": {
        "nextCursor": "abcdef",
		"totalItems": 2
    }
}`

const testResults2 = `{
    "data": [
		{"id": "6d497ca9-d5ef-495c-9074-a7021afc42c1", "name": "test2"}
    ],
	"pagination": {
        "nextCursor": null,
		"totalItems": 2
    }

}`

const testError = `{
	"headers": {
		"normalizedNames": {},
		"lazyUpdate": null
	},
	"status": 400,
	"statusText": "OK",
	"url": "https://euce1-test.sentinelone.net/web/api/v2.1/threats?cursor=test",
	"ok": false,
	"name": "HttpErrorResponse",
	"message": "Http failure response for https://euce1-test.sentinelone.net/web/api/v2.1/threats?cursor=test: 400 OK",
	"error": {
		"errors": [
			{
				"code": 4000010,
				"detail": "Invalid cursor value received",
				"title": "Validation JSONError"
			}
		]
	}
}`

func TestClient_GetJSONResponse(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/test",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, testData), nil
		})

	req, err := c.NewRequest("GET", "v2.1/test", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = c.GetJSONResponse(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClient_GetJSONResponse_Error(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/test",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, testError), nil
		})

	req, err := c.NewRequest("GET", "v2.1/test", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = c.GetJSONResponse(req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestClient_GetJSONItems(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/list",
		func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.RawQuery, "cursor=abcdef") {
				return httpmock.NewStringResponse(200, testResults2), nil
			}
			return httpmock.NewStringResponse(200, testResults), nil
		})

	req, err := c.NewRequest("GET", "v2.1/list", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = c.GetJSONItems(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
