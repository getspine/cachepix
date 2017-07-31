#!/usr/bin/env bash

set -ueo pipefail

cd /go/src/github.com/ssalevan/photocache

if [ "${PCACHE_REBUILD-false}" == true ]; then
  glide install
  go-wrapper install
fi

go-wrapper run "${@}"
