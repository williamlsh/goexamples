#!/usr/bin/env bash

docker run \
    -d \
    --rm \
    --name ipfs \
    -e IPFS_PROFILE=server \
    -v "$(pwd)"/staging:/export \
    -v "$(pwd)"/data:/data/ipfs \
    -p 4001:4001 \
    -p 4001:4001/udp \
    -p 127.0.0.1:8080:8080 \
    -p 127.0.0.1:5001:5001 \
    ipfs/go-ipfs:latest
