#!/usr/bin/env bash

VERSION_TO_TEST="$1"

docker pull monogradient/go-tracker-api:"${VERSION_TO_TEST}"

docker run -p