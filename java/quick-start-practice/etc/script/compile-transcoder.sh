#!/bin/bash


SHELL=/bin/bash

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=gitlab.alipay-inc.com,code.alipay.com

# build transcoder plugins
if [[ -n "${PLUGIN_TRANSCODER}" ]]; then
  coders=(${PLUGIN_TRANSCODER//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/transcoders
  for name in "${coders[@]}"; do
    export PLUGIN_TARGET=${name}
    export PLUGIN_TRANSCODER_OUTPUT=${PLUGIN_TRANSCODER_PREFIX}-${PLUGIN_TARGET}.so
    # check BUILD_OPTS
    if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
      build_opts="GOOS=${PLUGIN_OS} GOARCH=${PLUGIN_ARCH}"
      export BUILD_OPTS=${build_opts}
      echo "compiling transcoder ${name} for ${PLUGIN_OS} ${PLUGIN_ARCH} ..."
    else
      echo "compiling transcoder ${name} for linux $(dpkg --print-architecture) ..."
    fi
    make compile-transcoder
  done
fi
