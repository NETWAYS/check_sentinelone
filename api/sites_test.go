package api_test

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestClient_GetSites(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/sites",
		func(req *http.Request) (*http.Response, error) {
			data, err := os.ReadFile("testdata/sites.json")
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			return httpmock.NewBytesResponse(200, data), nil
		})

	sites, err := c.GetSites(url.Values{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sites) != 1 {
		t.Fatalf("expected sites to be of len 1")
	}
}

func TestClient_GetSites_Integration(t *testing.T) {
	c := envClient(t)

	sites, err := c.GetSites(url.Values{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if sites == nil {
		t.Fatalf("expected sites not being nil")
	}
}
