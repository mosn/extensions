#!/bin/bash


SHELL=/bin/bash

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=gitlab.alipay-inc.com,code.alipay.com

# build trace plugins
if [[ -n "${PLUGIN_TRACE}" ]]; then
  tracers=(${PLUGIN_TRACE//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/traces
  for name in "${tracers[@]}"; do
    export PLUGIN_TARGET=${name}
    export PLUGIN_TRACE_OUTPUT=${PLUGIN_TRACE_PREFIX}-${PLUGIN_TARGET}.so
    if [[ -n "${PLUGIN_GIT_VERSION}" ]]; then
      export PLUGIN_TRACE_OUTPUT=${PLUGIN_TRACE_PREFIX}-${PLUGIN_TARGET}-${PLUGIN_GIT_VERSION}.so
    fi
    # check BUILD_OPTS
    if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
      build_opts="GOOS=${PLUGIN_OS} GOARCH=${PLUGIN_ARCH}"
      export BUILD_OPTS=${build_opts}
      echo "compiling trace ${name} for ${PLUGIN_OS} ${PLUGIN_ARCH} ..."
    else
      echo "compiling trace ${name} for linux $(dpkg --print-architecture) ..."
    fi
    make compile-trace
  done
fi
