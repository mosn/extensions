#!/bin/bash

URL=http://127.0.0.1:2088/brn/pdf/0001/000016
#EXTREF=$ date +%s%N
EXTREF=$(date +%Y%m%d%H%M%S)$RANDOM
echo  $URL
echo $EXTREF

BODY="$(cat req.txt)"

echo $BODY
echo
echo "--------START CURL---------"
echo

curl --location --request POST "$URL" \
--header 'AreaCode: 0000' \
--header 'X-BOLE-SourceSysKey: 92b8a321739c4b539d654e9f739178a1' \
--header 'OrigSender: ESB002' \
--header 'VersionId: 0001' \
--header 'CtrlBits: 10000000' \
--header 'TraceId: esb001f4a8fe69615a474091e0a768fd15c2af' \
--header 'SpanId: esb001bc2128f41f054817b711f475453d51c3' \
--header 'Accept-Encoding: gzip,deflate' \
--header 'Content-Type: text/plain' \
--data-raw "$BODY"

echo
echo "--------END-------------"