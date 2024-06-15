#!/bin/bash


set -e

pushd $(dirname $0)/..

. scripts/client/common.sh

docker build \
    --platform linux/amd64 \
    -t $client_builder_image_name \
    -f scripts/client/Dockerfile \
    .

docker run \
    --rm \
    --platform linux/amd64 \
    -v $(pwd)/internal/client/dist:/src/client/dist:rw \
    -v $(pwd)/internal/client/src/generated:/src/client/src/generated:rw \
    $client_builder_image_name \
    npm run build
