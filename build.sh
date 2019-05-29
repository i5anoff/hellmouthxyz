#!/bin/bash
go version

GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/renderingobjectanimation/*.go
chmod +x main
