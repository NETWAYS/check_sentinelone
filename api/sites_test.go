package api_test

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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
			assert.NoError(t, err)

			return httpmock.NewBytesResponse(200, data), nil
		})

	sites, err := c.GetSites(url.Values{})
	assert.NoError(t, err)
	assert.Len(t, sites, 1)
}

func TestClient_GetSites_Integration(t *testing.T) {
	c := envClient(t)

	sites, err := c.GetSites(url.Values{})
	assert.NoError(t, err)
	assert.NotNil(t, sites)
}
