#!/bin/bash

# package transcoder plugins
if [[ -f "/go/src/${PLUGIN_PROJECT_NAME}/build/sidecar/binary/mosn" ]]; then
  PLUGIN_ANT_ZIP_OUTPUT=mosn.zip
  if [[ -f "/go/src/${PLUGIN_PROJECT_NAME}/etc/bundle/${PLUGIN_ANT_ZIP_OUTPUT}" ]]; then
    rm -rf "/go/src/${PLUGIN_PROJECT_NAME}/etc/bundle/${PLUGIN_ANT_ZIP_OUTPUT}"
  fi
  cd "/go/src/${PLUGIN_PROJECT_NAME}/build/sidecar/binary/"
  echo "packaging mosn..."
  zip -r "${PLUGIN_ANT_ZIP_OUTPUT}" .
  mv "/go/src/${PLUGIN_PROJECT_NAME}/build/sidecar/binary/${PLUGIN_ANT_ZIP_OUTPUT}" \
    "/go/src/${PLUGIN_PROJECT_NAME}/etc/bundle/${PLUGIN_ANT_ZIP_OUTPUT}"
fi
