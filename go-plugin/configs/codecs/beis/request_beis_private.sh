#!/bin/bash

IP=127.0.0.1
PORT=3045
echo
echo  $IP $PORT

echo
echo $(date +%Y-%m-%d,%H:%M:%s)

BODY="$(cat beis_req_private.txt)"

echo $BODY
echo
echo "--------START invoke---------"
echo

(sleep 1;echo "$BODY";sleep 300;) | nc $IP $PORT
echo
echo "--------END-------------"

#
############################################################
