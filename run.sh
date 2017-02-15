#!/usr/bin/env bash

echo "Browse to http://localhost:8000/ to test"

docker run -it -p 8000:8000 sverrirab/kube-test-container

