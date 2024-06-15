#!/bin/bash


set -e

pushd $(dirname $0)/..

scripts/validate-openapi.sh
scripts/build-client.sh

date_tag=$(date +%Y-%m-%d)
local_output_dir=$(pwd)/dist/$date_tag
mkdir -p $local_output_dir

cp bookmarklet.js $local_output_dir/
if [ -f cmd/cick-playlister/credentials.json ]; then
    cp cmd/cick-playlister/credentials.json $local_output_dir/
fi

image_name="captaincoordinates/cick-playlister-builder"

docker build \
    -t $image_name \
    -f scripts/api/Dockerfile \
    --build-arg date_tag=${date_tag} \
    .

docker run \
    --rm \
    -v $(pwd)/dist:/src/dist:rw \
    $image_name
