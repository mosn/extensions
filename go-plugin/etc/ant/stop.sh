#!/bin/bash

# kill container if running
sidecar=$(docker ps -a -q -f name=mosn-container)
if [[ -n "$sidecar" ]]; then
  echo "mosn-container is running and terminating..."
  docker stop mosn-container >/dev/null
  docker rm -f mosn-container >/dev/null
  echo "terminated ok"
else
  echo "no mosn-container is running"
fi
