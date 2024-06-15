#!/bin/bash


set -e

pushd $(dirname $0)/../cmd/cick-playlister
go build cick-playlister.go
