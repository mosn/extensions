#! /bin/bash

BASE_IMAGE=zonghaishang/delve:1.6.1
PROJECT_NAME=github.com/mosn/extensions/go-plugin
TARGET=${plugin}

sidecar=$(docker ps -a -q -f name=mosn-container)
if [[ -n "$sidecar" ]]; then
  docker stop mosn-container >/dev/null
  docker rm -f mosn-container >/dev/null
  echo "stop existed sidecar success."
fi

EXPORT_PORTS="-p 2345:2345 -p 2045:2045 -p 2046:2046"

docker run -u admin \
  -e PLUGIN_PROJECT_NAME="${PROJECT_NAME}" \
  -e PLUGIN_TARGET="${TARGET}" \
  -v $(go env GOPATH)/src/${PROJECT_NAME}:/go/src/${PROJECT_NAME} \
  -itd --name mosn-container --env-file env_conf ${EXPORT_PORTS} \
  -w /go/src/${PROJECT_NAME} \
  ${BASE_IMAGE} /go/src/${PROJECT_NAME}/etc/script/run.sh ${type} ${TARGET}

echo "start mosn-container container."
echo "run 'docker exec -it mosn-container /bin/bash' command enter mosn container."
