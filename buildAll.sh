#!/usr/bin/env bash

pushd $(dirname "$0")

set -ex

GOOS=linux   GOARCH=386   go build -o ./build/sda_linux_386
GOOS=windows GOARCH=386   go build -o ./build/sda_win_386
GOOS=darwin  GOARCH=amd64 go build -o ./build/sda_osx_amd64
GOOS=linux   GOARCH=amd64 go build -o ./build/sda_linux_amd64
GOOS=windows GOARCH=amd64 go build -o ./build/sda_win_amd64
