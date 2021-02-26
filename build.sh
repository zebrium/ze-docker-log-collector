#!/bin/sh
#
# Please note you should not run this script directly to build
# docker image. This script is run automatically inside docker
# container when we build image.
#
# To build container image, use the command below:
#
# docker build --no-cache -t <repo>/ze-docker-log-collector:<tag> -f Dockerfile .
#
# Please note code for building is from github.com/zebrium/ze-docker-log-collector/zebrium,
# not this repo. So you must commit and push changes before you make a build.
#

set -e
apk add --update go git mercurial build-base ca-certificates
mkdir -p /go/src/github.com/gliderlabs
cp -r /src /go/src/github.com/gliderlabs/logspout
cd /go/src/github.com/gliderlabs/logspout
export GOPATH=/go
go get
go build -ldflags "-X main.Version=$1" -o /bin/logspout
apk del go git mercurial build-base
rm -rf /go
rm -rf /var/cache/apk/*

# backwards compatibility
ln -fs /tmp/docker.sock /var/run/docker.sock
