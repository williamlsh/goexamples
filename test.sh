#!/usr/bin/env bash

curl -vv http://localhost:8080 \
  --data "hello world at $(date)"
