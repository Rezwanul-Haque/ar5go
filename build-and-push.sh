#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <tag>"
  exit 0
fi

go build
echo "Docker build with tag: $1"
docker build -t <docker-registry-path>/clean:"$1" .
docker push <docker-registry-path>/clean:"$1"
