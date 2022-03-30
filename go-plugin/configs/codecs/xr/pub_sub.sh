#!/bin/bash

echo "pub service CIMT000070@xr"
curl -X POST -d '{"protocolType": "xr", "providerMetaInfo": { "appName": "xr-provider","properties": {"application": "xr-provider","port": "9999" }},	"serviceName": "CIMT000070@xr"}' localhost:13330/services/publish

sleep 2

echo "sub service CIMT000070@xr"
curl -X POST -d '{"protocolType":"xr","serviceName":"CIMT000070@xr"}' localhost:13330/services/subscribe