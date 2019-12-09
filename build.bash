#!/bin/bash

if [[ -e app ]]; then
    rm app
fi

if [[ -e go-ldflags ]]; then
    rm go-ldflags
fi

# go build -ldflags="-X main.Version=v1.0.0" -o app

go build -v -ldflags="-X 'main.Version=v1.0.0' -X 'github.com/williamlsh/goexamples/build.Time=$(date)' -X 'github.com/williamlsh/goexamples/build.User=$(id -u -n)'" -o app
