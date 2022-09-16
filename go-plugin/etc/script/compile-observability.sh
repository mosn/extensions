#!/bin/bash

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=gitlab.alipay-inc.com,code.alipay.com

build_opts=""
if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
  build_opts="GOOS=${PLUGIN_OS} GOARCH=${PLUGIN_ARCH}"
  echo "compiling observability ${PLUGIN_TARGET} for ${PLUGIN_OS} ${PLUGIN_ARCH} ..."
else
  echo "compiling observability ${PLUGIN_TARGET} for linux $(dpkg --print-architecture) ..."
fi

export BUILD_OPTS=${build_opts}
export PLUGIN_OBSERVABILITY_OUTPUT=${PLUGIN_OBSERVABILITY_PREFIX}-${PLUGIN_TARGET}-${COMMIT}.so

# build trace 
if [[ -n "${PLUGIN_TRACE}" ]]; then
  bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-trace.sh
fi

if [[ -n "${PLUGIN_TRACE}" ]]; then
  trace=(${PLUGIN_TRACE//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/observability/trace 
  meta=
  for name in "${trace[@]}"; do
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/observability/trace
    echo "cp  /go/src/${PLUGIN_PROJECT_NAME}/build/trace/${name} /go/src/${PLUGIN_PROJECT_NAME}/build/observability/trace/"
    cp -r /go/src/${PLUGIN_PROJECT_NAME}/build/trace/${name} \
      /go/src/${PLUGIN_PROJECT_NAME}/build/observability/trace/
    # append ,
    if [[ -n ${meta} ]]; then
      meta+=","
    fi
    meta+="\"${name}\""
  done
fi

