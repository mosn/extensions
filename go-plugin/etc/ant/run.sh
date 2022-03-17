#! /bin/bash

mosn="/go/src/${PLUGIN_PROJECT_NAME}/build/sidecar/binary/mosn"
MOSN_PREFIX=/home/admin/mosn

mkdir -p $MOSN_PREFIX/conf \
  $MOSN_PREFIX/logs \
  $MOSN_PREFIX/bin \
  $MOSN_PREFIX/base_conf/clusters \
  $MOSN_PREFIX/conf/clusters \
  $MOSN_PREFIX/base_conf/certs

cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/process_checker.sh /home/admin/mosn/bin/process_checker.sh
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/update_checker.sh /home/admin/mosn/bin/update_checker.sh
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/zclean.sh /home/admin/mosn/bin/zclean.sh
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/zclean_crontab.sh /home/admin/mosn/bin/zclean_crontab.sh
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/gen-cert.sh /home/admin/mosn/base_conf/certs/gen-cert.sh
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/prestop.sh /home/admin/mosn/bin/prestop.sh

cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/export_node_port.py /home/admin/mosn/bin/export_node_port.py
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/modify_iptables.sh /home/admin/mosn/bin/modify_iptables.sh
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/iptables.hijack /home/admin/mosn/bin/iptables.hijack
cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/clean_iptables.sh /home/admin/mosn/bin/clean_iptables.sh

cp "${mosn}" /home/admin/mosn/bin/mosn

chmod +x /home/admin/mosn/bin/mosn
chown -R admin:admin /home/admin

cp /go/src/${SIDECAR_PROJECT_NAME}/configs/mosn_config.json $MOSN_PREFIX/base_conf/mosn_config.json
# COPY Basic configs
cp -r /go/src/${SIDECAR_PROJECT_NAME}/configs/routers $MOSN_PREFIX/base_conf/routers
cp -r /go/src/${SIDECAR_PROJECT_NAME}/configs/listeners $MOSN_PREFIX/base_conf/listeners
cp -r /go/src/${SIDECAR_PROJECT_NAME}/configs/certs $MOSN_PREFIX/base_conf/certs

echo "sidecar->  ${mosn}"

debug=${SIDECAR_DLV_DEBUG}
if [[ -n "$debug" && "$debug" == "true" ]]; then
  echo "running mode: debug"
  dlv --listen=0.0.0.0:2345 --headless=true --api-version=2 --accept-multiclient --allow-non-terminal-interactive exec /home/admin/mosn/bin/mosn -- start -c /home/admin/mosn/conf/mosn_config.json -b /home/admin/mosn/base_conf/mosn_config.json -n "sidecar~$RequestedIP~$POD_NAME.$POD_NAMESPACE~$POD_NAMESPACE.$DOMAINNAME" -s "$Sigma_Site" -l /home/admin/logs/mosn/default.log
else
  /home/admin/mosn/bin/mosn start -c /home/admin/mosn/conf/mosn_config.json -b /home/admin/mosn/base_conf/mosn_config.json -n "sidecar~$RequestedIP~$POD_NAME.$POD_NAMESPACE~$POD_NAMESPACE.$DOMAINNAME" -s "$Sigma_Site" -l /home/admin/logs/mosn/default.log
fi
