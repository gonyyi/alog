#!/bin/sh

go test -bench=. -json > ./tmp/bench.txt
