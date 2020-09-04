package api_test

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClient_GetThreats(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://euce1-test.sentinelone.net/web/api/v2.1/threats",
		func(req *http.Request) (*http.Response, error) {
			data, err := ioutil.ReadFile("testdata/threats.json")
			assert.NoError(t, err)

			return httpmock.NewBytesResponse(200, data), nil
		})

	threats, err := c.GetThreats()
	assert.NoError(t, err)
	assert.Len(t, threats, 1)
}

func TestClient_GetThreats_Integration(t *testing.T) {
	c := envClient(t)

	threats, err := c.GetThreats()
	assert.NoError(t, err)
	assert.NotNil(t, threats)
}
