#!/usr/bin/env bash

cd ./helm/charts/

helm package ../kube-test-container

helm repo index --url https://raw.githubusercontent.com/sverrirab/kube-test-container/master/helm/charts/ .
