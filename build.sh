#!/usr/bin/env bash

GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags kube-test-container ./cmd/kube-test-container.go

IMAGE_NAME=sverrirab/kube-test-container

docker build -t sverrirab/kube-test-container .

if [ -n "$1" ] ; then
    echo "Tagging with $1"
    docker tag sverrirab/kube-test-container sverrirab/kube-test-container:$1
fi
