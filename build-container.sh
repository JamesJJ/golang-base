#!/bin/sh

which gofmt &>/dev/null && find . -type f -path .git -prune -o -name '*.go' -exec gofmt -w {} \;

#export DOCKER_BUILDKIT=1

docker build $@ .
