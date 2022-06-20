#!/bin/bash

# please change ehttp-provider to your service identity
export SERVICE="ehttp-provider"
# please change userInfo to your request url
export REQUEST_URL=hello
export REQUEST_PORT=3045

REQUEST_COMMAND="curl -v -H \"X-SERVICE: ${SERVICE}\" -H \"Content-Type: application/json\" localhost:${REQUEST_PORT}/${REQUEST_URL}"

echo "${REQUEST_COMMAND}"

curl -v -H "X-SERVICE: ${SERVICE}" -H "Content-Type: application/json" localhost:${REQUEST_PORT}/${REQUEST_URL}
