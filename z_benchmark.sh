#!/bin/sh
mkdir -p ./tmp

if [ -z "$1" ] ; then
  # if no args given, run bench on screen
  go test -bench=.
else
  go test -bench=. > ./tmp/$1
  echo "Saved to ./tmp/$1"
fi
