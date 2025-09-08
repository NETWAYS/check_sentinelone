package api_test

import (
	"github.com/NETWAYS/check_sentinelone/api"
	"github.com/jarcoal/httpmock"
	"net/http"
	"os"
	"testing"
)

func envClient(t *testing.T) *api.Client {
	url := os.Getenv("SENTINELONE_URL")
	token := os.Getenv("SENTINELONE_TOKEN")

	if url == "" || token == "" {
		t.Skip("SENTINELONE_URL and SENTINELONE_TOKEN must be set!")
	}

	return api.NewClient(url, token)
}

func testClient() (*api.Client, func()) {
	httpmock.Activate()

	return api.NewClient("https://euce1-test.sentinelone.net", "test"), func() {
		httpmock.DeactivateAndReset()
	}
}

func TestClient(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/test",
		func(req *http.Request) (*http.Response, error) {
			body := map[string]any{
				"data": map[string]string{"test": "test"},
			}
			return httpmock.NewJsonResponse(200, body)
		})

	req, err := c.NewRequest("GET", "v2.1/test", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	res, err := c.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if 200 != res.StatusCode {
		t.Fatalf("expected %v, got %v", 200, res.StatusCode)
	}

	res.Body.Close()
}
