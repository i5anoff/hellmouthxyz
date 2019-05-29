#!/bin/bash
go version

GO111MODULE=on GOOS=linux go build -o main ./cmd/renderingobjectanimation/betterjson/main.go
chmod +x main
