# Defeway Toolbox

[![Build Status](https://travis-ci.org/crabtree/defeway-toolbox.svg?branch=master)](https://travis-ci.org/crabtree/defeway-toolbox)
[![Maintainability](https://api.codeclimate.com/v1/badges/25cd6143e39b5d2c5caa/maintainability)](https://codeclimate.com/github/crabtree/defeway-toolbox/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/crabtree/defeway-toolbox)](https://goreportcard.com/report/github.com/crabtree/defeway-toolbox)

## Build defeway-download binary

```
go build -o defewaydownload ./cmd/download
```

## Use defeway-download binary

Usage of `defewaydownload` binary:

- `-addr value` - IP address of the DVR
- `-port int` - port of the DVR (default 60001)
- `-username string` - username for the DVR (default "admin")
- `-password string` - password for the DVR (default empty)
- `-output string` - path to the downloads directory
- `-overwrite` - overwrite existing files
- `-concurrent int` - the number of concurrent workers (default 1)
- `-chan value` - channel id, you can specify multiple channels
- `-date value` - date in format YYYY-MM-DD (eg. 2019-01-01)
- `-start value` - recordings strat time
- `-end value` - recordings end time
- `-type value` - recording type, you can specify multiple types
- `-timeout timespan` - the timeout parameter for the HTTP client (default 5s)
- `-tls-skip-verify` - skip TLS verification
- `-no-keep-alives` - do not keep connections alive

## Build defeway-scan binary

```
go build -o defewayscan ./cmd/scan
```

## Use defeway-scan binary:

Usage of `defewayscan` binary:

- `-addr value` - IP address from which the scanner should start its job
- `-mask value` - network mask (eg. 255.255.255.0)
- `-port value` - the port of the DVR to scan, you can specify multiple ports
- `-username string` - username for the DVR (default "admin")
- `-password string` - password for the DVR (default empty)
- `-logdir string` - path to the logs directory
- `-concurrent int` - the number of concurrent workers (default 1)
- `-timeout timespan` - the timeout parameter for the HTTP client (default 5s)
- `-tls-skip-verify` - skip TLS verification
