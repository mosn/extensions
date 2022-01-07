#!/bin/bash

# package stream filter plugins
if [[ -n "${PLUGIN_STREAM_FILTER}" ]]; then
  filters=(${PLUGIN_STREAM_FILTER//,/ })
  for name in "${filters[@]}"; do
    PLUGIN_FILTER_ZIP_OUTPUT=${name}.zip
    rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/target/stream_filters/${PLUGIN_FILTER_ZIP_OUTPUT}
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/target/stream_filters/
    if [ -d "/go/src/${PLUGIN_PROJECT_NAME}/build/stream_filters/${name}/" ]; then
      cd /go/src/${PLUGIN_PROJECT_NAME}/build/stream_filters/
      echo "packaging filter ${name}..."
      zip -r /go/src/${PLUGIN_PROJECT_NAME}/build/target/stream_filters/${PLUGIN_FILTER_ZIP_OUTPUT} ${name} \
        -x "stream_filters/*" -x "transcoders/*" -x "mosn_config.json"
    fi
  done
fi
