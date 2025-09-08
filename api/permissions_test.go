package api_test

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestClient_CheckViewer(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/users/viewer-auth-check",
		func(req *http.Request) (*http.Response, error) {
			body := map[string]any{
				"data": map[string]any{"success": true},
			}
			return httpmock.NewJsonResponse(200, body)
		})

	ok, err := c.CheckViewer()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !ok {
		t.Fatalf("expected true, got false")
	}
}

func TestClient_CheckViewer_Integration(t *testing.T) {
	c := envClient(t)

	ok, err := c.CheckViewer()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !ok {
		t.Fatalf("expected true, got false")
	}
}
