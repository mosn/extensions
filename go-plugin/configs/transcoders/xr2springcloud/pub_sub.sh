#!/bin/bash

echo "pub service http.server@springcloud"
curl -X POST -d '{"protocolType": "springcloud", "providerMetaInfo": { "appName": "springcloud-provider","properties": {"application": "springcloud-provider","port": "18080" }},"serviceName": "http.server@springcloud"}' localhost:13330/services/publish

sleep 2

echo "sub service http.server@springcloud"
curl -X POST -d '{"protocolType":"springcloud","serviceName":"http.server@springcloud"}' localhost:13330/services/subscribe
