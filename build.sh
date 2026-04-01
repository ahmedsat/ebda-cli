#!/usr/bin/env bash

set -xe

export CGO_ENABLED=0
export GOOS=windows

go build -tags=release -o ~/Downloads/ebda-cli.exe .
