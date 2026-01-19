#!/bin/bash
set CGO_ENABLED=0
set GOOS=linux
go mod tidy
go build -o game_srv main.go