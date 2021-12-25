#!/usr/bin/env bash

: "${KO_DOCKER_REPO:?"You must set environment variable 'KO_DOCKER_REPO'"}"

export GO111MODULE=on
export GOFLAGS=-mod=vendor

cat << EOF | ko resolve -Bf -
images:
- ko://github.com/scribe-security/ci-demo-app/test/test_images/http
EOF
