#!/usr/bin/env bash
set -ex
export WEBCIGPORT=8081
export WEBCIG_TEMPLATE_DIR=${GOPATH}/src/github.com/deorbit/webcig/server/templates/
export WEBCIG_STATIC_DIR=${GOPATH}/src/github.com/deorbit/webcig/server/static/
go generate github.com/deorbit/webcig/webcigd
go install github.com/deorbit/webcig/webcigd
exec webcigd