package api

import (
	"encoding/json"
	"net/url"
	"time"
)

type Threat struct {
	AgentRealtimeInfo AgentRealtimeInfo `json:"agentRealtimeInfo"`
	ThreatInfo        ThreatInfo        `json:"threatInfo"`
}

type AgentRealtimeInfo struct {
	AccountID         string `json:"accountId"`
	AccountName       string `json:"accountName"`
	SiteID            string `json:"siteId"`
	SiteName          string `json:"siteName"`
	GroupID           string `json:"groupId"`
	GroupName         string `json:"groupName"`
	AgentComputerName string `json:"agentComputerName"`
	AgentDomain       string `json:"AgentDomain"`
}

type ThreatInfo struct {
	ThreatName                  string    `json:"threatName"`
	Classification              string    `json:"classification"`
	ClassificationSource        string    `json:"classificationSource"`
	CreatedAt                   time.Time `json:"createdAt"`
	Engines                     []string  `json:"engines"`
	AnalystVerdict              string    `json:"analystVerdict"`
	AnalystVerdictDescription   string    `json:"analystVerdictDescription"`
	IncidentStatus              string    `json:"incidentStatus"`
	IncidentStatusDescription   string    `json:"incidentStatusDescription"`
	MitigationStatus            string    `json:"mitigationStatus"`
	MitigationStatusDescription string    `json:"mitigationStatusDescription"`
}

func (c *Client) GetThreats(values url.Values, computername string) (threats []*Threat, err error) {
	// nolint: noctx
	req, err := c.NewRequest("GET", "v2.1/threats?"+values.Encode(), nil)
	if err != nil {
		return
	}

	res, err := c.GetJSONItems(req)
	if err != nil {
		return
	}

	for _, item := range res {
		t := &Threat{}

		err = json.Unmarshal(item, t)
		if err != nil {
			return
		}
		if computername != "" && computername != t.AgentRealtimeInfo.AgentComputerName {
			continue
		}
		threats = append(threats, t)
	}

	return
}
