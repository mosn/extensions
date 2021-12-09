#! /bin/bash

PLUGIN_TYPE=$1
PLUGIN_TARGET=$2

# handle alias plugin type
if [[ ${PLUGIN_TYPE} == "codec" ]];then
  PLUGIN_TYPE="codecs"
fi

if [[ ${PLUGIN_TYPE} == "trans" ]];then
  PLUGIN_TYPE="transcoders"
fi

if [[ ${PLUGIN_TYPE} == "sf" ]];then
  PLUGIN_TYPE="transcoders"
fi

mosn="/go/src/${PLUGIN_PROJECT_NAME}/build/sidecar/binary/mosn"
SIDECAR_CONF="/go/src/${PLUGIN_PROJECT_NAME}/build/${PLUGIN_TYPE}/${PLUGIN_TARGET}/mosn_config.json"

echo "----> ${SIDECAR_CONF}"

mkdir /home/admin/mosn/bin
mkdir -p /home/admin/logs

cp "${mosn}" /home/admin/mosn/bin/mosn
cp "${SIDECAR_CONF}" /home/admin/mosn/bin/mosn_config.json

chmod +x /home/admin/mosn/bin/mosn
chown -R admin:admin /home/admin

echo "sidecar->  ${mosn}"
echo "conf-> ${SIDECAR_CONF}"

dlv --listen=0.0.0.0:2345 --headless=true --api-version=2 --accept-multiclient --allow-non-terminal-interactive exec /home/admin/mosn/bin/mosn -- start -c /home/admin/mosn/bin/mosn_config.json
