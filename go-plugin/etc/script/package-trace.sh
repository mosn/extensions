#!/bin/bash

# build trace plugins
if [[ -n "${PLUGIN_TRACE}" ]]; then
  traces=(${PLUGIN_TRACE//,/ })
  for name in "${traces[@]}"; do
    if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
      export PLUGIN_TARGET=${name}
      export PLUGIN_TRACE_OUTPUT=${PLUGIN_TRACE_PREFIX}-${PLUGIN_TARGET}.so
      bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-trace.sh
    elif [[ "${PLUGIN_BUILD_PLATFORM}" == "Darwin" && "${PLUGIN_BUILD_PLATFORM_ARCH}" == "arm64" ]]; then
      # apple m1 chip compile plugin(amd64)
      export PLUGIN_TARGET=${name}
      export PLUGIN_TRACE_OUTPUT=${PLUGIN_TRACE_PREFIX}-${PLUGIN_TARGET}.so

      export PLUGIN_OS="linux"
      export PLUGIN_ARCH="amd64"
      bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-trace.sh
    fi
  done
fi

# package trace plugins
if [[ -n "${PLUGIN_TRACE}" ]]; then
  traces=(${PLUGIN_TRACE//,/ })
  for name in "${traces[@]}"; do
    PLUGIN_TRACE_ZIP_OUTPUT=${name}.zip
    rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/target/trace/${PLUGIN_TRACE_ZIP_OUTPUT}
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/target/trace/
    if [ -d "/go/src/${PLUGIN_PROJECT_NAME}/build/trace/${name}/" ]; then
      cd /go/src/${PLUGIN_PROJECT_NAME}/build/trace/
      echo "packaging trace ${name}..."
      zip -r /go/src/${PLUGIN_PROJECT_NAME}/build/target/trace/${PLUGIN_TRACE_ZIP_OUTPUT} ${name} 
    fi
  done
fi
