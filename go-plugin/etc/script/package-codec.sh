#!/bin/bash

# build codec plugins
if [[ -n "${PLUGIN_CODEC}" ]]; then
  coders=(${PLUGIN_CODEC//,/ })
  for name in "${coders[@]}"; do
    if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
      export PLUGIN_TARGET=${name}
      export PLUGIN_CODEC_OUTPUT=${PLUGIN_CODEC_PREFIX}-${PLUGIN_TARGET}.so
      bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-codec.sh
    elif [[ "${PLUGIN_BUILD_PLATFORM}" == "Darwin" && "${PLUGIN_BUILD_PLATFORM_ARCH}" == "arm64" ]]; then
      # apple m1 chip compile plugin(amd64)
      export PLUGIN_TARGET=${name}
      export PLUGIN_CODEC_OUTPUT=${PLUGIN_CODEC_PREFIX}-${PLUGIN_TARGET}.so

      export PLUGIN_OS="linux"
      export PLUGIN_ARCH="amd64"
      bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-codec.sh
    fi
  done
fi

pkg_version=
if [[ -f "/go/src/${PLUGIN_PROJECT_NAME}/VERSION.txt" ]]; then
  pkg_version=$(cat /go/src/${PLUGIN_PROJECT_NAME}/VERSION.txt)
fi

# package codec plugins
if [[ -n "${PLUGIN_CODEC}" ]]; then
  coders=(${PLUGIN_CODEC//,/ })
  for name in "${coders[@]}"; do
    PLUGIN_CODEC_ZIP_OUTPUT=${name}.zip
    if [[ -n "${pkg_version}" ]]; then
      PLUGIN_CODEC_ZIP_OUTPUT=${name}-${pkg_version}.zip
    fi
    rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/target/codecs/${PLUGIN_CODEC_ZIP_OUTPUT}
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/target/codecs/
    if [ -d "/go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${name}/" ]; then
      cd /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/
      echo "packaging codec ${name}..."
      zip -r /go/src/${PLUGIN_PROJECT_NAME}/build/target/codecs/${PLUGIN_CODEC_ZIP_OUTPUT} ${name} \
        -x "stream_filters/*" -x "transcoders/*" -x "mosn_config.json"
    fi
  done
fi
