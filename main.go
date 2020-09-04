package main

import (
	"github.com/NETWAYS/go-check"
)

const readme = `Check for threats on the SentinelOne Cloud service.

You need to provide the URL of your instance and an authentication token, which is user specific.
It is recommended to create a new user with "Viewer" permissions only.

Threats will be listed until their incident state has been resolved, or with the
--ignore-in-progress flag, is no longer "unresolved". Mitigated threats appear as warning.

https://github.com/NETWAYS/check_sentinelone

Copyright (c) 2020 NETWAYS GmbH <info@netways.de>
Copyright (c) 2020 Markus Frosch <markus.frosch@netways.de

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).`

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
