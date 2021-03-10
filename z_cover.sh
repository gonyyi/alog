#!/bin/sh
mkdir -p ./tmp
go test -coverprofile ./tmp/coverage.out && go tool cover -html=./tmp/coverage.out
