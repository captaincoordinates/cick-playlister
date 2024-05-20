#!/bin/bash


set -e

pushd $(dirname $0)/..

scripts/validate-openapi.sh

pushd cmd/cick-playlister
go build
