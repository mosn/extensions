#!/bin/bash

echo "pub service com.alipay.sofa.ms.service.EchoService@dubbo"
curl -X POST -d '{"protocolType": "dubbo", "providerMetaInfo": { "appName": "dubbo-provider","properties": {"application": "dubbo-provider","port": "20880" }},	"serviceName": "com.alipay.sofa.ms.service.EchoService@dubbo"}' localhost:13330/services/publish

sleep 2

echo "sub service com.alipay.sofa.ms.service.EchoService@dubbo"
curl -X POST -d '{"protocolType":"dubbo","serviceName":"com.alipay.sofa.ms.service.EchoService@dubbo"}' localhost:13330/services/subscribe
