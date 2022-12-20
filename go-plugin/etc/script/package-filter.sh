#!/bin/bash

# build stream filter plugins
if [[ -n "${PLUGIN_STREAM_FILTER}" ]]; then
  if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
    bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-filter.sh
  elif [[ "${PLUGIN_BUILD_PLATFORM}" == "Darwin" && "${PLUGIN_BUILD_PLATFORM_ARCH}" == "arm64" ]]; then
    # apple m1 chip compile plugin(amd64)
    export PLUGIN_OS="linux"
    export PLUGIN_ARCH="amd64"
    bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-filter.sh
  fi
fi

# package stream filter plugins
if [[ -n "${PLUGIN_STREAM_FILTER}" ]]; then
  filters=(${PLUGIN_STREAM_FILTER//,/ })
  for name in "${filters[@]}"; do
    PLUGIN_FILTER_ZIP_OUTPUT=${name}.zip
    rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/target/stream_filters/${PLUGIN_FILTER_ZIP_OUTPUT}
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/target/stream_filters/
    if [ -d "/go/src/${PLUGIN_PROJECT_NAME}/build/output/stream_filters/${name}/" ]; then
      cd /go/src/${PLUGIN_PROJECT_NAME}/build/output/stream_filters/
      echo "packaging filter ${name}..."
      zip -r /go/src/${PLUGIN_PROJECT_NAME}/build/target/stream_filters/${PLUGIN_FILTER_ZIP_OUTPUT} ${name} \
        -x "stream_filters/*" -x "transcoders/*" -x "mosn_config.json"
    fi
  done
fi
