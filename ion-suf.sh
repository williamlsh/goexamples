#!/usr/bin/env bash

sudo docker run \
    -d \
    --rm \
    --name sfu \
    -p 7000:7000 \
    -p 5000-5020:5000-5020/udp \
    pionwebrtc/ion-sfu:latest-jsonrpc
