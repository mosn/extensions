#!/bin/bash

URL=http://127.0.0.1:10088/reservations/echo?message=hello
#EXTREF=$ date +%s%N
echo  $URL

echo
echo "--------START CURL---------"
echo

curl --location --request GET "$URL" \
--header 'Content-Type: text/plain' \
--header 'X-Target-App: reservation-client' \

echo
echo "--------END-------------"
