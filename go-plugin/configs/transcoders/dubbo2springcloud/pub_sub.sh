#!/bin/bash

echo "pub service reservation-service@springcloud"
curl -X POST -d '{"protocolType": "springcloud", "providerMetaInfo": { "appName": "springcloud-provider","properties": {"application": "springcloud-provider","port": "18080" }},"serviceName": "reservation-service@springcloud"}' localhost:13330/services/publish

sleep 2

echo "sub service reservation-service@springcloud"
curl -X POST -d '{"protocolType":"springcloud","serviceName":"reservation-service@springcloud"}' localhost:13330/services/subscribe
