#!/bin/bash

echo "pub service ECIF120000300003@cd"
curl -X POST -d '{"protocolType": "cd", "providerMetaInfo": { "appName": "cd-provider","properties": {"application": "cd-provider","port": "10150" }},	"serviceName": "ECIF120000300003@cd"}' localhost:13330/services/publish

sleep 2

echo "sub service ECIF120000300003@cd"
curl -X POST -d '{"protocolType":"cd","serviceName":"ECIF120000300003@cd"}' localhost:13330/services/subscribe