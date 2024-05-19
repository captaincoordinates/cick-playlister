#!/bin/bash


set -e

pushd $(dirname $0)/..

docker compose build openapi
docker compose run openapi
