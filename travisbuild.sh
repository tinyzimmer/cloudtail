#!/usr/bin/env bash

set -e
set -o pipefail
set -x

gox -os="linux windows" -arch="amd64" -output="./build/{{.Dir}}-{{.OS}}-{{.Arch}}"
gox -os="darwin" -arch="amd64" -output="./build/{{.Dir}}-macOS-{{.Arch}}"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stbuild/cloudtail .
docker build . -t cloudtail
echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
docker tag cloudtail tinyzimmer/cloudtail
docker push tinyzimmer/cloudtail
