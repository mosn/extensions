#!/bin/bash

mosn="/go/src/${PLUGIN_PROJECT_NAME}/build/sidecar/binary/mosn"
MOSN_PREFIX=/home/admin/mosn

mkdir -p $MOSN_PREFIX/conf \
  $MOSN_PREFIX/logs \
  $MOSN_PREFIX/bin \
  $MOSN_PREFIX/base_conf/clusters \
  $MOSN_PREFIX/conf/clusters \
  $MOSN_PREFIX/base_conf/certs

if [[ -d "/go/src/${SIDECAR_PROJECT_NAME}" ]]; then

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/process_checker.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/process_checker.sh /home/admin/mosn/bin/process_checker.sh
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/update_checker.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/update_checker.sh /home/admin/mosn/bin/update_checker.sh
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/zclean.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/zclean.sh /home/admin/mosn/bin/zclean.sh
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/zclean_crontab.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/zclean_crontab.sh /home/admin/mosn/bin/zclean_crontab.sh
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/gen-cert.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/gen-cert.sh /home/admin/mosn/base_conf/certs/gen-cert.sh
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/prestop.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/prestop.sh /home/admin/mosn/bin/prestop.sh
  fi

  if [[ -f "go/src/${SIDECAR_PROJECT_NAME}/etc/script/export_node_port.py" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/export_node_port.py /home/admin/mosn/bin/export_node_port.py
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/modify_iptables.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/modify_iptables.sh /home/admin/mosn/bin/modify_iptables.sh
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/iptables.hijack" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/iptables.hijack /home/admin/mosn/bin/iptables.hijack
  fi
  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/etc/script/clean_iptables.sh" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/etc/script/clean_iptables.sh /home/admin/mosn/bin/clean_iptables.sh
  fi

  # COPY Basic configs

  if [[ -d "/go/src/${SIDECAR_PROJECT_NAME}/configs/routers" ]]; then
    cp -r /go/src/${SIDECAR_PROJECT_NAME}/configs/routers $MOSN_PREFIX/base_conf/routers
  fi

  if [[ -d "/go/src/${SIDECAR_PROJECT_NAME}/configs/listeners" ]]; then
    cp -r /go/src/${SIDECAR_PROJECT_NAME}/configs/listeners $MOSN_PREFIX/base_conf/listeners
  fi

  if [[ -d "/go/src/${SIDECAR_PROJECT_NAME}/configs/certs" ]]; then
    cp -r /go/src/${SIDECAR_PROJECT_NAME}/configs/certs $MOSN_PREFIX/base_conf/certs
  fi

  if [[ -f "/go/src/${SIDECAR_PROJECT_NAME}/configs/mosn_config.json" ]]; then
    cp /go/src/${SIDECAR_PROJECT_NAME}/configs/mosn_config.json $MOSN_PREFIX/base_conf/mosn_config.json
  else
    # no source available, remove source level dependency
    cp -r /go/src/${PLUGIN_PROJECT_NAME}/configs/internal/base_conf/routers $MOSN_PREFIX/base_conf/routers
    cp -r /go/src/${PLUGIN_PROJECT_NAME}/configs/internal/base_conf/listeners $MOSN_PREFIX/base_conf/listeners
    cp /go/src/${PLUGIN_PROJECT_NAME}/configs/internal/base_conf/mosn_config.json $MOSN_PREFIX/base_conf/mosn_config.json
  fi
fi

cp "${mosn}" /home/admin/mosn/bin/mosn

chmod +x /home/admin/mosn/bin/mosn
chown -R admin:admin /home/admin
chmod +x /go/src/${PLUGIN_PROJECT_NAME}/etc/ant/*.sh
chmod +x /go/src/${PLUGIN_PROJECT_NAME}/etc/script/*.sh

echo "sidecar->  ${mosn}"

debug=${SIDECAR_DLV_DEBUG}
if [[ -n "$debug" && "$debug" == "true" ]]; then
  echo "running mode: debug, arch: $(dpkg --print-architecture)"
  dlv --listen=0.0.0.0:2345 --continue --headless=true --api-version=2 --accept-multiclient --allow-non-terminal-interactive exec /home/admin/mosn/bin/mosn -- start -c /home/admin/mosn/conf/mosn_config.json -b /home/admin/mosn/base_conf/mosn_config.json -n "sidecar~$RequestedIP~$POD_NAME.$POD_NAMESPACE~$POD_NAMESPACE.$DOMAINNAME" -s "$Sigma_Site" -l /home/admin/logs/mosn/default.log
else
  /home/admin/mosn/bin/mosn start -c /home/admin/mosn/conf/mosn_config.json -b /home/admin/mosn/base_conf/mosn_config.json -n "sidecar~$RequestedIP~$POD_NAME.$POD_NAMESPACE~$POD_NAMESPACE.$DOMAINNAME" -s "$Sigma_Site" -l /home/admin/logs/mosn/default.log
fi
