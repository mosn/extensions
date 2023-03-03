#!/bin/bash

BASE_IMAGE=zonghaishang/delve:v1.7.3

DEBUG_PORTS="-p 2345:2345"
LISTENER_PORTS="-p 11001:11001 -p 12220:12220 -p 12200:12200 -p 30880:30880 -p 30800:30800 -p 10088:10088 -p 10080:10080 -p 34904:34904 -p 15001:15001 -p 15006:15006 -p 3399:3399 -p 13399:13399"
EXPORT_PORTS="-p 34901:34901 -p 13330:13330 -p 2045:2045 -p 2046:2046 -p 16379:16379 -p 9529:9529 -p 9530:9530"

# biz port export
BIZ_PORTS=" -p 13088:13088 -p 13080:13080"

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
os_name=$(uname)
if [[ "$os_name" == "Linux" ]]; then
    export PUB_BOLT_LOCAL_IP=$(ip -f inet address | grep inet | grep -v docker | grep -v 127.0.0.1 | head -n1 | awk '{print $2}' | awk -F/ '{print $1}')
  else
    # default for mac os
    export PUB_BOLT_LOCAL_IP=$(ipconfig getifaddr en0)
fi
echo "host address: ${PUB_BOLT_LOCAL_IP} ->  ${PROJECT_NAME}"

# create mapping logs directory
if [[ -d ${FULL_PROJECT_NAME}/logs ]]; then
   rm -rf ${FULL_PROJECT_NAME}/logs
fi
mkdir -p ${FULL_PROJECT_NAME}/logs

docker run ${DOCKER_BUILD_OPTS} \
  -u admin --privileged \
  -e PLUGIN_PROJECT_NAME="${PROJECT_NAME}" \
  -e DYNAMIC_CONF_PATH=/go/src/${PROJECT_NAME}/build/codecs \
  -e SIDECAR_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} \
  -e SIDECAR_DLV_DEBUG="${DEBUG_MODE}" \
  -v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
  -v ${FULL_PROJECT_NAME}/logs:/home/admin/logs \
  -v $(go env GOPATH)/src/${SIDECAR_GITLAB_PROJECT_NAME}:/go/src/${SIDECAR_GITLAB_PROJECT_NAME} \
  -d --name mosn-container --env-file "${FULL_PROJECT_NAME}"/etc/ant/env_conf ${MAPPING_PORTS} \
  -w /go/src/${PROJECT_NAME} \
  ${BASE_IMAGE} /go/src/${PROJECT_NAME}/etc/ant/run.sh "$@"

echo "start mosn-container container success."
echo "run 'docker exec -it mosn-container /bin/bash' command enter mosn container."
