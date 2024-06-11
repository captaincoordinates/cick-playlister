#!/bin/bash


set -e

pushd $(dirname $0)/..

pushd internal/client
npm install
npm run watch
