# Defeway Toolbox

[![Maintainability](https://api.codeclimate.com/v1/badges/25cd6143e39b5d2c5caa/maintainability)](https://codeclimate.com/github/crabtree/defeway-toolbox/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/crabtree/defeway-toolbox)](https://goreportcard.com/report/github.com/crabtree/defeway-toolbox)

## Build defeway-download binary

```
go build -o defewaydownload ./cmd/download
```

## Use defeway-download binary

Usage of `defewaydownload` binary:

- `-addr value` - IP address of the DVR
- `-chan value` - channel id, you can specify multiple channels, optional when `-file` specified
- `-concurrent int` - the number of concurrent workers (default 1)
- `-date value` - date in format YYYY-MM-DD (eg. 2019-01-01)
- `-end value` - recordings end time
- `-file string` - path to the XML file with a list of recordings to download
- `-no-keep-alives` - do not keep connections alive
- `-output string` - path to the downloads directory
- `-overwrite` - overwrite existing files
- `-password string` - password for the DVR (default empty)
- `-port int` - port of the DVR (default 60001)
- `-preview` - limit the length of the downloads to about 1 minute
- `-start value` - recordings strat time
- `-timeout timespan` - the timeout parameter for the HTTP client (default 5s)
- `-tls-skip-verify` - skip TLS verification
- `-type value` - recording type, you can specify multiple types, optional when `-file` specified
- `-username string` - username for the DVR (default "admin")

## Build defeway-scan binary

```
go build -o defewayscan ./cmd/scan
```

## Use defeway-scan binary:

Usage of `defewayscan` binary:

- `-addr value` - IP address from which the scanner should start its job
- `-concurrent int` - the number of concurrent workers (default 1)
- `-logdir string` - path to the logs directory
- `-mask value` - network mask (eg. 255.255.255.0)
- `-port value` - the port of the DVR to scan, you can specify multiple ports
- `-password string` - password for the DVR (default empty)
- `-timeout timespan` - the timeout parameter for the HTTP client (default 5s)
- `-tls-skip-verify` - skip TLS verification
- `-username string` - username for the DVR (default "admin")
