#! /bin/bash

BASE_IMAGE=zonghaishang/delve:1.6.1
PROJECT_NAME=github.com/mosn/extensions/go-plugin
SIDECAR_GITLAB_PROJECT_NAME=gitlab.alipay-inc.com/ant-mesh/mosn

DEBUG_PORTS="-p 2345:2345"
LISTENER_PORTS="-p 12220:12220 -p 12200:12200 -p 30880:30880 -p 30800:30800 -p 10088:10088 -p 10080:10080 -p 34904:34904 -p 15001:15001 -p 15006:15006 -p 3399:3399 -p 13399:13399"
EXPORT_PORTS="-p 2045:2045 -p 2046:2046 -p 13330:13330 -p 16379:16379 -p 9529:9529 -p 9530:9530 -p 34901:34901"

# biz port export
BIZ_PORTS="-p 13088:13088 -p 13080:13080"

MAPPING_PORTS="${DEBUG_PORTS} ${LISTENER_PORTS} ${EXPORT_PORTS} ${BIZ_PORTS}"

sidecar=$(docker ps -a -q -f name=mosn-container)
if [[ -n "$sidecar" ]]; then
  echo "found mosn-container is running already and terminating..."
  docker stop mosn-container >/dev/null
  docker rm -f mosn-container >/dev/null
  rm -rf $(go env GOPATH)/src/${PROJECT_NAME}/logs
  echo "terminated ok"
fi

# export local ip for mosn
export PUB_BOLT_LOCAL_IP=$(ifconfig -a | grep inet | grep -v 127.0.0.1 | grep -v inet6 | grep -v "inet 0" | awk '{print $2}' | tr -d "addr:")
echo "host address: ${PUB_BOLT_LOCAL_IP}"

docker run -u admin \
  -e PLUGIN_PROJECT_NAME="${PROJECT_NAME}" \
  -e DYNAMIC_CONF_PATH=/go/src/${PROJECT_NAME}/build/codecs \
  -e SIDECAR_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} \
  -v $(go env GOPATH)/src/${PROJECT_NAME}:/go/src/${PROJECT_NAME} \
  -v $(go env GOPATH)/src/${PROJECT_NAME}/logs:/home/admin/logs \
  -v $(go env GOPATH)/src/${SIDECAR_GITLAB_PROJECT_NAME}:/go/src/${SIDECAR_GITLAB_PROJECT_NAME} \
  -itd --name mosn-container --env-file $(go env GOPATH)/src/${PROJECT_NAME}/etc/ant/env_conf ${MAPPING_PORTS} \
  -w /go/src/${PROJECT_NAME} \
  ${BASE_IMAGE} /go/src/${PROJECT_NAME}/etc/ant/run.sh "$@"

echo "start mosn-container container success."
echo "run 'docker exec -it mosn-container /bin/bash' command enter mosn container."
