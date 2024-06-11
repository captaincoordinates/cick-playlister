#!/bin/bash


set -e

pushd $(dirname $0)/..

scripts/validate-openapi.sh
scripts/build-client.sh

output_dir=$(pwd)/dist/$(date +%Y-%m-%d)
mkdir -p $output_dir
output_path=$output_dir/click-playlister.exe

cp bookmarklet.js $output_dir/

pushd cmd/cick-playlister
GOOS=windows GOARCH=amd64 go build -o $output_path cick-playlister.go
if [ -f credentials.json ]; then
    cp credentials.json $output_dir/
fi
