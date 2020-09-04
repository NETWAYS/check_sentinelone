package main

import (
	"errors"
	"fmt"
	"github.com/NETWAYS/check_sentinelone/api"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/pflag"
	"net/url"
	"os"
)

type Config struct {
	ManagementURL    string
	AuthToken        string
	IgnoreInProgress bool
}

func BuildConfigFlags(fs *pflag.FlagSet) (config *Config) {
	config = &Config{}

	fs.StringVarP(&config.ManagementURL, "url", "H", "", "Management URL (env:SENTINELONE_URL)")
	fs.StringVarP(&config.AuthToken, "token", "T", "", "API AuthToken (env:SENTINELONE_TOKEN)")

	fs.BoolVar(&config.IgnoreInProgress, "ignore-in-progress", false,
		"Ignore threats, where the incident status is in-progress")

	return
}

func (c *Config) SetFromEnv() {
	if c.ManagementURL == "" {
		c.ManagementURL = os.Getenv("SENTINELONE_URL")
	}

	if c.AuthToken == "" {
		c.AuthToken = os.Getenv("SENTINELONE_TOKEN")
	}

	return
}

func (c *Config) Validate() (err error) {
	if c.ManagementURL == "" || c.AuthToken == "" {
		err = errors.New("url and token are required")
		return
	}

	return
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

	threats, err := client.GetThreats(values)
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

	for index, list := range byLocation {
		output += fmt.Sprintf("\n## %s\n\n", index)

		for _, threat := range list {
			total++

			if threat.ThreatInfo.MitigationStatus == "not_mitigated" {
				notMitigated++
			}

			output += fmt.Sprintf("[%s] %s: (%s) %s (%s)\n",
				threat.ThreatInfo.CreatedAt.Local().Format("2006-01-02 15:04 MST"),
				threat.AgentRealtimeInfo.AgentComputerName,
				threat.ThreatInfo.Classification,
				threat.ThreatInfo.ThreatName,
				threat.ThreatInfo.MitigationStatusDescription,
			)
		}
	}

	// Add summary on top
	output = fmt.Sprintf("%d threats found, %d not mitigated\n", total, notMitigated) + output

	if notMitigated > 0 {
		rc = check.Critical
	} else if total > 0 {
		rc = check.Warning
	}

	return
}
