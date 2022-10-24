#bin/bash

rm -f ./log.txt && nohup go run main.go  > log.txt 2>&1 &