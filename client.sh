#!/usr/bin/env bash

curl -s --request "POST" \
    --header "Content-Type: application/json" \
    --data '{"inches": 1}' \
    http://localhost:8080/twirp/haberdasher.Haberdasher/MakeHat |
    jq

echo 'inches:1' |
    protoc --encode haberdasher.Size ./pb/service.proto |
    curl -s --request POST \
        --header "Content-Type: application/protobuf" \
        --data-binary @- \
        http://localhost:8080/twirp/haberdasher.Haberdasher/MakeHat |
    protoc --decode haberdasher.Hat ./pb/service.proto
