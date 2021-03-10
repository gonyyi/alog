#!/bin/sh
mkdir -p ./tmp
go test -bench=. -json > ./tmp/bench.txt
