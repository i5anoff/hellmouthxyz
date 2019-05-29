#!/bin/bash
go version

GO111MODULE=on go build -o main ./cmd/renderingobjectanimation/*.go
chmod +x main
