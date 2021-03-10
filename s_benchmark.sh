#!/bin/sh
mkdir -p ./tmp

if [ -z "$1" ] ; then
  # if no args given, run bench on screen
  go test -bench=.
else
  go test -bench=. $1 > ./tmp/bench.txt
  echo "Saved to ./tmp/bench.txt"
fi
