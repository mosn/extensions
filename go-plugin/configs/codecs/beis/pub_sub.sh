#!/bin/bash

echo "pub service BRNC140030280001@beis"
curl -X POST -d '{"protocolType": "beis", "providerMetaInfo": { "appName": "beis-provider","properties": {"application": "beis-provider","port": "7766" }},	"serviceName": "BRNC140030280001@beis"}' localhost:13330/services/publish

sleep 2

echo "sub service BRNC140030280001@beis"
curl -X POST -d '{"protocolType":"beis","serviceName":"BRNC140030280001@beis"}' localhost:13330/services/subscribe