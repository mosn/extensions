#!/bin/bash

URL=http://127.0.0.1:10088/reservations/echo
#EXTREF=$ date +%s%N
echo  $URL
echo hello

echo
echo "--------START CURL---------"
echo

curl --location --request POST "$URL" \
--header 'Content-Type: text/plain' \
--header 'X-Target-App: reservation-client' \
--data-raw '"hello"'

echo
echo "--------END-------------"
