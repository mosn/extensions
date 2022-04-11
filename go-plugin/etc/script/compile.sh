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

build_opts="CGO_ENABLED=1"

if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
  export GOOS=${PLUGIN_OS}
  export GOARCH=${PLUGIN_ARCH}

  build_opts="${build_opts} GOOS=${GOOS} GOARCH=${GOARCH}"
  echo "compiling mosn for ${PLUGIN_OS} ${PLUGIN_ARCH} ..."
else
  echo "compiling mosn for linux $(dpkg --print-architecture) ..."
fi

export CGO_ENABLED=1

echo "${build_opts} go build -o mosn ${SIDECAR_PROJECT}/cmd/mosn/main"

go build -mod=readonly -gcflags "all=-N -l" \
  -ldflags "-B 0x$(head -c20 /dev/urandom | od -An -tx1 | tr -d ' \n') -X main.Version=${MAJOR_VERSION} -X main.GitVersion=${GIT_VERSION}" \
  -v -o mosn "${SIDECAR_PROJECT}/cmd/mosn/main"

if [ -f mosn ]; then
  md5sum -b mosn | cut -d' ' -f1 >mosn-${MAJOR_VERSION}-${GIT_VERSION}.md5
  mv mosn-${MAJOR_VERSION}-${GIT_VERSION}.md5 "/go/src/${PLUGIN_PROJECT}/build/sidecar/binary/mosn-${MAJOR_VERSION}-${GIT_VERSION}.md5"
  mv mosn "/go/src/${PLUGIN_PROJECT}/build/sidecar/binary/mosn"
fi
