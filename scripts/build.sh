#!/bin/bash


set -e

pushd $(dirname $0)/..

scripts/validate-openapi.sh
scripts/build-client.sh

pushd cmd/cick-playlister
go build
