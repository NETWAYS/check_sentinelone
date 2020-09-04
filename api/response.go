package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ResponseBody struct {
	Data       json.RawMessage `json:"data"`
	Pagination *Pagination     `json:"pagination"`
	// Seems to be only used with errors
	Name string `json:"name"`
	// Seems to be only used with errors
	Message string           `json:"message"`
	Error   *JSONErrorObject `json:"error"`
	Errors  []JSONError      `json:"errors"`
}

type Pagination struct {
	NextCursor string `json:"nextCursor"`
	TotalItems int    `json:"totalItems"`
}

type SuccessMessage struct {
	Success bool `json:"success"`
}

type JSONErrorObject struct {
	Errors []JSONError `json:"errors"`
}

type JSONError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (c *Client) GetJSONResponse(req *http.Request) (data *ResponseBody, err error) {
	res, err := c.Do(req)
	if err != nil {
		return
	}

	// read response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("could not retrieve response body: %w", err)
		return
	}

	res.Body.Close()

	// Try to parse json
	data = &ResponseBody{}
	err = json.Unmarshal(body, data)

	if res.StatusCode != 200 {
		errInfo := ""

		if data.Error != nil {
			for _, e := range data.Error.Errors {
				errInfo += fmt.Sprintf(" - %s: %s", e.Title, e.Detail)
			}
		}

		for _, e := range data.Errors {
			errInfo += fmt.Sprintf(" - %s: %s", e.Title, e.Detail)
		}

		if err == nil {
			log.WithFields(log.Fields{
				"status": res.StatusCode,
				"errors": data.Errors,
			}).Debug("HTTP returned non-ok result")
		} else {
			log.WithFields(log.Fields{
				"status": res.StatusCode,
				"body":   string(body),
			}).Debug("HTTP returned non-ok result without JSON info")
		}

		err = fmt.Errorf("HTTP request returned non-ok status %s%s", res.Status, errInfo)

		return
	}

	if err != nil {
		err = fmt.Errorf("could not decode JSON from body: %w", err)
		return
	}

	return
}

func (c *Client) GetJSONItems(request *http.Request) (items []json.RawMessage, err error) {
	var (
		ctx        = request.Context()
		nextCursor string
		response   *ResponseBody
	)

	for {
		r := request.Clone(ctx)
		if nextCursor != "" {
			if r.URL.RawQuery != "" {
				r.URL.RawQuery += "&"
			}

			r.URL.RawQuery += "cursor=" + url.QueryEscape(nextCursor)
		}

		response, err = c.GetJSONResponse(r)
		if err != nil {
			return
		}

		// retrieve items from response
		var dataItems []json.RawMessage

		err = json.Unmarshal(response.Data, &dataItems)
		if err != nil {
			// Fall back to adding response.Data to the item list
			// This is useful when data is not an array, but an object
			err = nil

			items = append(items, response.Data)
		} else {
			// append each item to overall list
			for _, item := range dataItems {
				items = append(items, item)
			}
		}

		// set nextCursor or break iteration when done
		if response.Pagination.NextCursor == "" {
			break
		} else if response.Pagination.NextCursor == nextCursor {
			err = fmt.Errorf("iteration error in pages, nextCursor is the same as before: %s", nextCursor)
			return
		} else {
			nextCursor = response.Pagination.NextCursor
		}
	}

	return
}
