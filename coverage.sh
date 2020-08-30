#!/usr/bin/env bash

go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
open coverage.html