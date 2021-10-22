#!/bin/sh

echo "Testing ..."
go test ./...

echo "Building ..."
[ ! -d bin ] && mkdir bin

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/defewaydownload ./cmd/download
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/defewaydownload.exe ./cmd/download

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/defewayscan ./cmd/scan
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/defewayscan.exe ./cmd/scan

echo "... DONE!"