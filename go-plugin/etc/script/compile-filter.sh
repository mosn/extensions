#!/bin/bash

# build stream filter plugins
if [[ -n "${PLUGIN_STREAM_FILTER}" ]]; then
  filters=(${PLUGIN_STREAM_FILTER//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/stream_filters
  for name in "${filters[@]}"; do
    export PLUGIN_TARGET=${name}
    export PLUGIN_STEAM_FILTER_OUTPUT=${PLUGIN_STEAM_FILTER_PREFIX}-${PLUGIN_TARGET}.so
    make compile-stream-filter
  done
fi
