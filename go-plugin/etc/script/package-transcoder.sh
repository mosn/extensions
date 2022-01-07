#!/bin/bash

# package transcoder plugins
if [[ -n "${PLUGIN_TRANSCODER}" ]]; then
  codecs=(${PLUGIN_TRANSCODER//,/ })
  for name in "${codecs[@]}"; do
    PLUGIN_TRANSCODER_ZIP_OUTPUT=${name}.zip
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
