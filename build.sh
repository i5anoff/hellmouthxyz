#!/bin/bash
go version

GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o main ./cmd/renderingobjectanimation/*.go
chmod +x main
