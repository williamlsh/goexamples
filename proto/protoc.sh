#!usr/bin/env bash

protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. greeter.proto