#!/bin/bash

# package trace plugins
if [[ -n "${PLUGIN_TRACE}" ]]; then
  if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
    bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-trace.sh
  elif [[ "${PLUGIN_BUILD_PLATFORM}" == "Darwin" && "${PLUGIN_BUILD_PLATFORM_ARCH}" == "arm64" ]]; then
    # apple m1 chip compile plugin(amd64)
    export PLUGIN_OS="linux"
    export PLUGIN_ARCH="amd64"
    bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-trace.sh
  fi
fi

pkg_version=
if [[ -f "/go/src/${PLUGIN_PROJECT_NAME}/VERSION.txt" ]]; then
  pkg_version=$(cat /go/src/${PLUGIN_PROJECT_NAME}/VERSION.txt)
fi

# package transcoder plugins
if [[ -n "${PLUGIN_TRACE}" ]]; then
  tracers=(${PLUGIN_TRACE//,/ })
  for name in "${tracers[@]}"; do
    PLUGIN_TRACE_ZIP_OUTPUT=${name}.zip
    if [[ -n "${pkg_version}" ]]; then
      PLUGIN_TRACE_ZIP_OUTPUT=${name}-${pkg_version}.zip
    fi
    rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/target/traces/${PLUGIN_TRACE_ZIP_OUTPUT}
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/target/traces/
    if [ -d "/go/src/${PLUGIN_PROJECT_NAME}/build/traces/${name}/" ]; then
      cd /go/src/${PLUGIN_PROJECT_NAME}/build/traces/
      echo "packaging trace ${name}..."
      zip -r /go/src/${PLUGIN_PROJECT_NAME}/build/target/traces/${PLUGIN_TRACE_ZIP_OUTPUT} ${name} \
        -x "stream_filters/*" -x "transcoders/*" -x "mosn_config.json"
    fi
  done
fi
