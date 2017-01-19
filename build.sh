#!/usr/bin/env bash

GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags kube-test-container ./cmd/kube-test-container.go

docker build -t sverrirab/kube-test-container:v1.0 .
