#!/bin/bash

echo "pub service NBT302800@bums"
curl -X POST -d '{"protocolType": "bums", "providerMetaInfo": { "appName": "bums-provider","properties": {"application": "bums-provider","port": "8999" }},	"serviceName": "NBT302800@bums"}' localhost:13330/services/publish

sleep 2

echo "sub service NBT302800@bums"
curl -X POST -d '{"protocolType":"bums","serviceName":"NBT302800@bums"}' localhost:13330/services/subscribe