#!/bin/bash

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=gitlab.alipay-inc.com,code.alipay.com

# build stream filter plugins
if [[ -n "${PLUGIN_TRACE}" ]]; then
  trace=(${PLUGIN_TRACE//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/trace
  for name in "${trace[@]}"; do
    export PLUGIN_TARGET=${name}
    export PLUGIN_TRACE_OUTPUT=${PLUGIN_TRACE_PREFIX}-${PLUGIN_TARGET}-${COMMIT}.so
    # check BUILD_OPTS
    if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
      build_opts="GOOS=${PLUGIN_OS} GOARCH=${PLUGIN_ARCH}"
      export BUILD_OPTS=${build_opts}
      echo "compiling filter ${name} for ${PLUGIN_OS} ${PLUGIN_ARCH} ..."
    else
      echo "compiling filter ${name} for linux $(dpkg --print-architecture) ..."
    fi
    make compile-trace
  done
fi
