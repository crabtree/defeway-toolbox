#!/bin/sh

echo "Testing ..."
go test ./...

echo "Building ..."
[ ! -d bin ] && mkdir bin

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/defewaydownload ./cmd/download
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/defewaydownload-amd64.exe ./cmd/download
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/defewaydownload-x86.exe ./cmd/download

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/defewayscan ./cmd/scan
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/defewayscan-amd64.exe ./cmd/scan
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/defewayscan-x86.exe ./cmd/scan

echo "... DONE!"