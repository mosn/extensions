#!/bin/bash

# build transcoder plugins
if [[ -n "${PLUGIN_TRANSCODER}" ]]; then
  coders=(${PLUGIN_TRANSCODER//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/transcoders
  for name in "${coders[@]}"; do
    export PLUGIN_TARGET=${name}
    export PLUGIN_TRANSCODER_OUTPUT=${PLUGIN_TRANSCODER_PREFIX}-${PLUGIN_TARGET}.so
    make compile-transcoder
  done
fi
