#!/bin/bash

SHELL=/bin/bash

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=gitlab.alipay-inc.com,code.alipay.com

export PLUGIN_PROJECT=${PROJECT_NAME}
export SIDECAR_PROJECT=${SIDECAR_PROJECT_NAME}

# update sidecar dependency.
go mod tidy
go mod download

MAJOR_VERSION=$(cat VERSION)
GIT_VERSION=$(git log -1 --pretty=format:%h)

rm -rf "/go/src/${PLUGIN_PROJECT}/build/sidecar/binary/"
mkdir -p "/go/src/${PLUGIN_PROJECT}/build/sidecar/binary/"

echo "go build -o mosn ${SIDECAR_PROJECT}/cmd/mosn/main"

CGO_ENABLED=1 go build -mod=readonly -gcflags "all=-N -l" \
  -ldflags "-B 0x$(head -c20 /dev/urandom | od -An -tx1 | tr -d ' \n') -X main.Version=${MAJOR_VERSION} -X main.GitVersion=${GIT_VERSION}" \
  -v -o mosn "${SIDECAR_PROJECT}/cmd/mosn/main"

mv mosn "/go/src/${PLUGIN_PROJECT}/build/sidecar/binary/mosn"
