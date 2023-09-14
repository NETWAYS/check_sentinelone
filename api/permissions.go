package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserInfo struct {
	ID       string            `json:"id"`
	IDType   string            `json:"idType"`
	APIHosts map[string]string `json:"apiHosts"`
}

func (c *Client) CheckViewer() (ok bool, err error) {
	// nolint: noctx
	req, err := c.NewRequest(http.MethodGet, "v2.1/users/viewer-auth-check", nil)
	if err != nil {
		return
	}

	resp, err := c.GetJSONResponse(req)
	if err != nil {
		return
	}

	message := &SuccessMessage{}

	err = json.Unmarshal(resp.Data, message)
	if err != nil {
		err = fmt.Errorf("could not decode SuccessMessage: %w", err)
		return
	}

	ok = message.Success

	return
}
