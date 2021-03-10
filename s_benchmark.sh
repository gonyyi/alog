#!/bin/sh
mkdir -p ./tmp

if [ -z "$1" ] ; then
  # if no args given, run bench on screen
  go test -bench=.
else
  if [ -z "$2" ] ; then
    go test -bench=. $1 > ./tmp/bench.json
    echo "Saved to ./tmp/bench.json"
  else
    go test -bench=. $1 > ./tmp/$2
    echo "Saved to ./tmp/$2"
  fi

fi
