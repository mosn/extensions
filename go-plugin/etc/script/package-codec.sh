#!/bin/bash

# package transcoder plugins
if [[ -n "${PLUGIN_CODEC}" ]]; then
  coders=(${PLUGIN_CODEC//,/ })
  for name in "${coders[@]}"; do
    PLUGIN_CODEC_ZIP_OUTPUT=${name}.zip
    rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/target/codecs/${PLUGIN_CODEC_ZIP_OUTPUT}
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/target/codecs/
    if [ -d "/go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${name}/" ]; then
      cd /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${name}
      echo "packaging codec ${name}..."
      zip -r /go/src/${PLUGIN_PROJECT_NAME}/build/target/codecs/${PLUGIN_CODEC_ZIP_OUTPUT} . \
        -x "stream_filters/*" -x "transcoders/*" -x "mosn_config.json"
    fi
  done
fi
