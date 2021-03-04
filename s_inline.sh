#!/bin/sh
go build -gcflags="-m -m" 2>&1 |grep "cannot" |sort > out.txt
# go build -gcflags="-m -m" 2>&1 | sort > out.txt && sublime out.txt

echo "\n-------\n" >> out.txt 

go build -gcflags="-m -m" 2>&1 |grep "can inline" |sort >> out.txt && sublime out.txt