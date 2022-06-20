#!/bin/bash

BASE_IMAGE=zonghaishang/delve:v1.7.3

DEBUG_PORTS="-p 2345:2345"
LISTENER_PORTS="-p  34901:34901"
EXPORT_PORTS="-p 13330:13330"

# biz port export
BIZ_PORTS=" -p 3045:3045 -p 3046:3046"

MAPPING_PORTS="${DEBUG_PORTS} ${LISTENER_PORTS} ${EXPORT_PORTS} ${BIZ_PORTS}"

sidecar=$(docker ps -a -q -f name=mosn-container)
if [[ -n "$sidecar" ]]; then
  echo
  echo "found mosn-container is running already and terminating..."
  docker stop mosn-container >/dev/null
  docker rm -f mosn-container >/dev/null
  rm -rf "${FULL_PROJECT_NAME}/logs"
  echo "terminated ok"
  echo
fi

DEBUG_MODE=${DLV_DEBUG}

chmod +x etc/ant/run.sh

# export local ip for mosn
export PUB_BOLT_LOCAL_IP=$(ipconfig getifaddr en0)
echo "host address: ${PUB_BOLT_LOCAL_IP} ->  ${PROJECT_NAME}"

docker run ${DOCKER_BUILD_OPTS} \
  -u admin --privileged \
  -e PLUGIN_PROJECT_NAME="${PROJECT_NAME}" \
  -e DYNAMIC_CONF_PATH=/go/src/${PROJECT_NAME}/build/codecs \
  -e SIDECAR_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} \
  -e SIDECAR_DLV_DEBUG="${DEBUG_MODE}" \
  -v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
  -v ${FULL_PROJECT_NAME}/logs:/home/admin/logs \
  -v $(go env GOPATH)/src/${SIDECAR_GITLAB_PROJECT_NAME}:/go/src/${SIDECAR_GITLAB_PROJECT_NAME} \
  -itd --name mosn-container --env-file "${FULL_PROJECT_NAME}"/etc/ant/env_conf ${MAPPING_PORTS} \
  -w /go/src/${PROJECT_NAME} \
  ${BASE_IMAGE} /go/src/${PROJECT_NAME}/etc/ant/run.sh "$@"

echo "start mosn-container container success."
echo "run 'docker exec -it mosn-container /bin/bash' command enter mosn container."
