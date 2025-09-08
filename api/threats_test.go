package api_test

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestClient_GetThreats(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/threats",
		func(req *http.Request) (*http.Response, error) {
			data, err := os.ReadFile("testdata/threats.json")
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			return httpmock.NewBytesResponse(200, data), nil
		})

	threats, err := c.GetThreats(url.Values{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(threats) != 1 {
		t.Fatalf("expected threats to be of len 1")
	}

}

func TestClient_GetThreats_Integration(t *testing.T) {
	c := envClient(t)

	threats, err := c.GetThreats(url.Values{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if threats == nil {
		t.Fatalf("expected threats not being nil")
	}
}
