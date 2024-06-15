#!/bin/bash


set -e

pushd $(dirname $0)/..

image_name="captaincoordinates/cick-playlister-openapi"

docker build \
    --platform linux/amd64 \
    -t $image_name \
    -f scripts/openapi/Dockerfile \
    .

docker run \
    --rm \
    --platform linux/amd64 \
    -v $(pwd)/internal/docs/openapi.yml:/openapi/openapi.yml:ro \
    $image_name
