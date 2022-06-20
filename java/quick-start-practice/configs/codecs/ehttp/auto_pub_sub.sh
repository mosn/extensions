#!/bin/bash

export SERVICE_ID="ehttp-provider@ehttp" # please change ehttp-provider to your service identity
export BACKEND_PORT=7755            # please change port 7755 to your java server port
export PROVIDER_APP=ehttp-provider

export MOCK_PUB_DATA="{\"protocolType\": \"ehttp\", \"providerMetaInfo\": { \"appName\": \"${PROVIDER_APP}\",\"properties\": {\"application\": \"${PROVIDER_APP}\",\"port\": \"${BACKEND_PORT}\" }},	\"serviceName\": \"${SERVICE_ID}\"}"

export MOCK_SUB_DATA="{\"protocolType\":\"ehttp\",\"serviceName\":\"${SERVICE_ID}\"}"

echo "publish service ${SERVICE_ID}"
echo "curl -d \"${MOCK_PUB_DATA}\" localhost:13330/services/publish"
curl -d "${MOCK_PUB_DATA}" localhost:13330/services/publish

sleep 2

echo
echo
echo "subscribe service ${SERVICE_ID}"
echo "curl -d \"${MOCK_SUB_DATA}\" localhost:13330/services/subscribe"
curl -d "${MOCK_SUB_DATA}" localhost:13330/services/subscribe
