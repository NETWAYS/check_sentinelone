package main

import (
	"github.com/NETWAYS/go-check"
)

const readme = `Check for threats on the SentinelOne Cloud service.

You need to provide the URL of your instance and an authentication token, which is user specific.
It is recommended to create a new user with "Viewer" permissions only.

Threats will be listed until their incident state has been resolved, or with the
--ignore-in-progress flag, is no longer "unresolved". Mitigated threats appear as warning.`

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()
	plugin.Name = "check_sentinelone"
	plugin.Readme = readme
	plugin.Version = buildVersion()
	plugin.Timeout = 30

	// Parse arguments
	config := BuildConfigFlags(plugin.FlagSet)
	plugin.ParseArguments()
	config.SetFromEnv()

	err := config.Validate()
	if err != nil {
		check.ExitError(err)
	}

	rc, output, err := config.Run()
	if err != nil {
		check.ExitError(err)
	}

	check.Exit(rc, output)
}
