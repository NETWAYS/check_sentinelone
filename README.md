check_sentinelone
=================

Check for threats on the SentinelOne Cloud service.

You need to provide the URL of your instance and an authentication token, which is user specific.
It is recommended to create a new user with "Viewer" permissions only.

Threats will be listed until their incident state has been resolved, or with the `--ignore-in-progress` flag, is no
longer "unresolved". Mitigated threats appear as warning.

## Usage

```
Arguments:
  -H, --url string             Management URL (e.g. https://your-site.sentinelone.net) (env:SENTINELONE_URL)
  -T, --token string           API AuthToken (env:SENTINELONE_TOKEN)
      --site string            Only list threats belonging to a named site
      --ignore-in-progress     Ignore threats, where the incident status is in-progress
      --computer-name string   Only list threats belonging to the specified computer name
  -t, --timeout int            Abort the check after n seconds (default 30)
  -d, --debug                  Enable debug mode
  -v, --verbose                Enable verbose mode
  -V, --version                Print version and exit
```

## Example

```
$ check_sentinelone --url https://your-site.sentinelone.net --token secret --site Customer

[CRITICAL] - site Customer - 13 threats found, 3 not mitigated
\_ [WARNING] fileserver: (Downloader) PDFCreator-1_9_4-setup.exe (Marked as benign)
\_ [WARNING] fileserver: (PUA) cdbxp_setup_4.5.7.6321.exe (Mitigated)
\_ [WARNING] fileserver: (Trojan) 2-1.exe (Mitigated)
\_ [WARNING] fileserver: (Trojan) 4-0.exe (Mitigated)
\_ [WARNING] fileserver: (Trojan) 7-0.exe (Mitigated)
\_ [WARNING] fileserver: (Trojan) 13-0.exe (Mitigated)
\_ [WARNING] fileserver: (Trojan) 1-0.exe (Mitigated)
\_ [WARNING] fileserver: (Trojan) 14-0.exe (Mitigated)
\_ [WARNING] fileserver: (Trojan) 12-0.exe (Mitigated)
\_ [CRITICAL] fileserver: (Adware) cdbxp_setup_4.5.8.7035.exe (Not mitigated)
\_ [WARNING] fileserver: (Adware) cdbxp_setup_4-{DFBDE0DF-DBEC-4437-A6D6-76CD670E9503}-v297222.exe (Mitigated)
\_ [CRITICAL] fileserver: (Adware) cdbxp_setup_4.5.8.7035.exe (Not mitigated)
\_ [CRITICAL] fileserver: (Adware) cdbxp_setup_4.5.8.7035.exe (Not mitigated)
| threats=13 threats_not_mitigated=3
```

## API Documentation

Full API documentation is available in `api-docs` under your Sentinel One dashboard.

This is only available for customers as fair as we know.

## License

Copyright (c) 2020 [NETWAYS GmbH](mailto:info@netways.de) \

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).
