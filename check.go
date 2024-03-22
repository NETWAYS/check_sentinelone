package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/NETWAYS/check_sentinelone/api"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/pflag"
)

type Config struct {
	ManagementURL    string
	AuthToken        string
	IgnoreInProgress bool
	SiteName         string
	ComputerName     string
}

func BuildConfigFlags(fs *pflag.FlagSet) (config *Config) {
	config = &Config{}

	fs.StringVarP(&config.ManagementURL, "url", "H", "",
		"Management URL (e.g. https://your-site.sentinelone.net) (env:SENTINELONE_URL)")
	fs.StringVarP(&config.AuthToken, "token", "T", "", "API AuthToken (env:SENTINELONE_TOKEN)")

	fs.StringVar(&config.SiteName, "site", "", "Only list threats belonging to a named site")

	fs.BoolVar(&config.IgnoreInProgress, "ignore-in-progress", false,
		"Ignore threats, where the incident status is in-progress")

	fs.StringVar(&config.ComputerName, "computer-name", "",
		"Only list threats belonging to the specified computer name")

	return
}

func (c *Config) SetFromEnv() {
	if c.ManagementURL == "" {
		c.ManagementURL = os.Getenv("SENTINELONE_URL")
	}

	if c.AuthToken == "" {
		c.AuthToken = os.Getenv("SENTINELONE_TOKEN")
	}
}

func (c *Config) Validate() error {
	if c.ManagementURL == "" || c.AuthToken == "" {
		return errors.New("url and token are required")
	}

	return nil
}

func (c *Config) Run() (rc int, output string, err error) {
	client := api.NewClient(c.ManagementURL, c.AuthToken)

	values := url.Values{}
	values.Set("sortOrder", "desc")

	if c.IgnoreInProgress {
		values.Set("incidentStatuses", "unresolved")
	} else {
		values.Set("resolved", "false")
	}

	if c.SiteName != "" {
		var siteID string

		siteID, err = lookupSiteID(client, c.SiteName)
		if err != nil {
			return
		}

		values.Set("siteIds", siteID)
	}

	if c.ComputerName != "" {
		values.Set("computerName__contains", c.ComputerName)
	}

	threats, err := client.GetThreats(values, c.ComputerName)
	if err != nil {
		return
	}

	var (
		total        int
		notMitigated int
	)

	byLocation := map[string][]*api.Threat{}

	for _, threat := range threats {
		index := fmt.Sprintf("%s / %s / %s",
			threat.AgentRealtimeInfo.AccountName,
			threat.AgentRealtimeInfo.SiteName,
			threat.AgentRealtimeInfo.GroupName,
		)

		if _, ok := byLocation[index]; !ok {
			byLocation[index] = []*api.Threat{}
		}

		byLocation[index] = append(byLocation[index], threat)
	}

	var sb strings.Builder

	for index, list := range byLocation {
		sb.WriteString(fmt.Sprintf("\n## %s\n\n", index))

		for _, threat := range list {
			var stateText string

			total++

			if threat.ThreatInfo.MitigationStatus == "not_mitigated" {
				notMitigated++

				stateText = "CRITICAL"
			} else {
				stateText = "WARNING"
			}

			sb.WriteString(fmt.Sprintf("[%s] [%s] %s: (%s) %s (%s)\n",
				// nolint: gosmopolitan
				threat.ThreatInfo.CreatedAt.Local().Format("2006-01-02 15:04 MST"),
				stateText,
				threat.AgentRealtimeInfo.AgentComputerName,
				threat.ThreatInfo.Classification,
				threat.ThreatInfo.ThreatName,
				threat.ThreatInfo.MitigationStatusDescription,
			))
		}
	}

	// Add summary on top.
	sb.WriteString(fmt.Sprintf("%d threats found, %d not mitigated\n", total, notMitigated) + sb.String())

	if c.SiteName != "" && c.ComputerName == "" {
		sb.WriteString(fmt.Sprintf("site %s - ", c.SiteName) + sb.String())
	} else if c.ComputerName != "" {
		sb.WriteString(fmt.Sprintf("Computer %s - ", c.ComputerName) + sb.String())
	}

	// Add perfdata.
	sb.WriteString("|")
	sb.WriteString(fmt.Sprintf(" threats=%d", total))
	sb.WriteString(fmt.Sprintf(" threats_not_mitigated=%d", notMitigated))
	output = sb.String()

	// determine final state.
	if notMitigated > 0 {
		rc = check.Critical
	} else if total > 0 {
		rc = check.Warning
	}

	return
}

func lookupSiteID(client *api.Client, name string) (id string, err error) {
	params := url.Values{}
	params.Set("name", name)
	params.Set("state", "active")

	sites, err := client.GetSites(params)
	if err != nil {
		return
	}

	switch len(sites) {
	case 0:
		err = fmt.Errorf("could not find a site named '%s'", name)
	case 1:
		id = sites[0].ID
	default:
		err = fmt.Errorf("more than one site matches '%s'", name)
	}

	return
}
