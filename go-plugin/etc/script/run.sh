#! /bin/bash

mosn="/go/src/${PLUGIN_PROJECT_NAME}/build/sidecar/binary/mosn"
SIDECAR_CONF="/go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/mosn_config.json"

mkdir /home/admin/bin
mkdir -p /home/admin/logs

cp "${mosn}" /home/admin/bin/mosn
cp "${SIDECAR_CONF}" /home/admin/bin/mosn_config.json

chmod +x /home/admin/bin/mosn
chown -R admin:admin /home/admin

echo "sidecar->  ${mosn}"
echo "conf-> ${SIDECAR_CONF}"

dlv --listen=0.0.0.0:2345 --headless=true --api-version=2 --accept-multiclient --allow-non-terminal-interactive exec /home/admin/bin/mosn -- start -c /home/admin/bin/mosn_config.json
