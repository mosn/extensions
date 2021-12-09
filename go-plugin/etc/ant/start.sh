#! /bin/bash

BASE_IMAGE=zonghaishang/delve:1.6.1
PROJECT_NAME=github.com/mosn/extensions/go-plugin
SIDECAR_GITLAB_PROJECT_NAME=gitlab.alipay-inc.com/ant-mesh/mosn

EXPORT_PORTS="-p 2045:2045 -p 2046:2046 -p 2345:2345 -p 13330:13330 -p 12200:12200 -p 12220:12220 -p 16379:16379 -p 9529:9529 -p 9530:9530 -p 34904:34904 -p 34901:34901"

# biz port export
OPTS=""

sidecar=$(docker ps -a -q -f name=mosn-container)
if [[ -n "$sidecar" ]]; then
  docker stop mosn-container >/dev/null
  docker rm -f mosn-container >/dev/null
  echo "stop existed sidecar success."
fi

# export local ip for mosn
export PUB_BOLT_LOCAL_IP=$(ifconfig -a | grep inet | grep -v 127.0.0.1 | grep -v inet6 | awk '{print $2}' | tr -d "addr:")
echo "host address: ${PUB_BOLT_LOCAL_IP}"

docker run -u admin \
  -e PLUGIN_PROJECT_NAME="${PROJECT_NAME}" \
  -e DYNAMIC_CONF_PATH=/go/src/${PROJECT_NAME}/build/codecs \
  -e SIDECAR_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} \
  -v $(go env GOPATH)/src/${PROJECT_NAME}:/go/src/${PROJECT_NAME} \
  -v $(go env GOPATH)/src/${SIDECAR_GITLAB_PROJECT_NAME}:/go/src/${SIDECAR_GITLAB_PROJECT_NAME} \
  -it --name mosn-container --env-file env_conf ${EXPORT_PORTS} \
  -w /go/src/${PROJECT_NAME} \
  ${BASE_IMAGE} /go/src/${PROJECT_NAME}/etc/ant/run.sh "$@"

echo "start mosn-container container."
echo "run 'docker exec -it mosn-container /bin/bash' command enter mosn container."
