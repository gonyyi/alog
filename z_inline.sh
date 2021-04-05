#!/bin/sh
mkdir -p ./tmp

go build -gcflags="-m -m" 2>&1 |grep "cannot" |sort > ./tmp/inline.txt

echo "\n-------\n" >> ./tmp/inline.txt 

go build -gcflags="-m -m" 2>&1 |grep "can inline" |sort >> ./tmp/inline.txt && sublime ./tmp/inline.txt
