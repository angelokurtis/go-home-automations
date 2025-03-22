#!/usr/bin/env bash

set -xe

# docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 -t kurtis/automations:latest --push .
