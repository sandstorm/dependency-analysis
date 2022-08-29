#!/usr/bin/env bash

pushd $(dirname "$0")

VERSION_FULL="$1"
VERSION_MAJOR=$(echo "$VERSION_FULL" | ack -o '^v\d+')

# validate version
echo "$VERSION_FULL" | ack '^v\d+\.\d+\.\d+$' > /dev/null
VERSION_EXITCODE=$?
if [ "$VERSION_EXITCODE" == "0" ]; then
    echo "creating release $VERSION_FULL (in short $VERSION_MAJOR)"
    read -p "Then press enter to continue."
    echo "Did you update the version in the version command?"
    read -p "Then press enter to continue."
else
    echo "invalid version format, should look like v1.2.3 and not like '$VERSION_FULL'"
    exit 1
fi

set -e

# build binaries
GOOS=linux   GOARCH=386   go build -o ./build/sda_linux_386
GOOS=windows GOARCH=386   go build -o ./build/sda_win_386
GOOS=darwin  GOARCH=amd64 go build -o ./build/sda_osx_amd64
GOOS=linux   GOARCH=amd64 go build -o ./build/sda_linux_amd64
GOOS=windows GOARCH=amd64 go build -o ./build/sda_win_amd64

# build Docker image
docker build -t sandstormmedia/dependency-analysis:latest .
docker tag sandstormmedia/dependency-analysis:latest sandstormmedia/dependency-analysis:v1

# tag version
git tag $VERSION_FULL
echo "Please execute: git push origin $VERSION_FULL"

# create release in Github
open https://github.com/sandstorm/dependency-analysis/releases/new

## TODO
#
# release Docker images
docker push sandstormmedia/dependency-analysis:latest
docker push sandstormmedia/dependency-analysis:v1
