package api

import (
	"encoding/json"
	"net/url"
)

type Site struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SiteType string `json:"siteType"`
}

type SiteResult struct {
	// Only license info
	AllSites json.RawMessage `json:"allSites"`
	Sites    []*Site
}

func (c *Client) GetSites(values url.Values) (data []*Site, err error) {
	req, err := c.NewRequest("GET", "v2.1/sites?state=active&"+values.Encode(), nil)
	if err != nil {
		return
	}

	res, err := c.GetJSONItems(req)
	if err != nil {
		return
	}

	for _, page := range res {
		p := &SiteResult{}

		err = json.Unmarshal(page, p)
		if err != nil {
			return
		}

		for _, site := range p.Sites {
			data = append(data, site)
		}
	}

	return
}
