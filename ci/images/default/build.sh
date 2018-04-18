#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

GOOS=linux GOARCH=amd64 go build -o "$DIR/boshua" $DIR/../../../cli/client/main.go

cd "$DIR"

docker build -t dpb587/bosh-compiled-releases-v2:latest .
