#!/bin/bash

# package transcoder plugins
if [[ -n "${PLUGIN_TRANSCODER}" ]]; then
  if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
    bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-transcoder.sh
  elif [[ "${PLUGIN_BUILD_PLATFORM}" == "Darwin" && "${PLUGIN_BUILD_PLATFORM_ARCH}" == "arm64" ]]; then
    # apple m1 chip compile plugin(amd64)
    export PLUGIN_OS="linux"
    export PLUGIN_ARCH="amd64"
    bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-transcoder.sh
  fi
fi

pkg_version=
if [[ -f "/go/src/${PLUGIN_PROJECT_NAME}/VERSION.txt" ]]; then
  pkg_version=$(cat /go/src/${PLUGIN_PROJECT_NAME}/VERSION.txt)
fi

# package transcoder plugins
if [[ -n "${PLUGIN_TRANSCODER}" ]]; then
  codecs=(${PLUGIN_TRANSCODER//,/ })
  for name in "${codecs[@]}"; do
    PLUGIN_TRANSCODER_ZIP_OUTPUT=${name}.zip
    if [[ -n "${pkg_version}" ]]; then
      PLUGIN_TRANSCODER_ZIP_OUTPUT=${name}-${pkg_version}.zip
    fi
    rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/target/transcoders/${PLUGIN_TRANSCODER_ZIP_OUTPUT}
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/target/transcoders/
    if [ -d "/go/src/${PLUGIN_PROJECT_NAME}/build/transcoders/${name}/" ]; then
      cd /go/src/${PLUGIN_PROJECT_NAME}/build/transcoders/
      echo "packaging transcoder ${name}..."
      zip -r /go/src/${PLUGIN_PROJECT_NAME}/build/target/transcoders/${PLUGIN_TRANSCODER_ZIP_OUTPUT} ${name} \
        -x "stream_filters/*" -x "transcoders/*" -x "mosn_config.json"
    fi
  done
fi
