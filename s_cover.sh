#!/bin/sh
go test -coverprofile ./_tmp/cover.out && go tool cover -html=./_tmp/cover.out