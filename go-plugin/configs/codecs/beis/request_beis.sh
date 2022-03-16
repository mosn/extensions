#!/bin/bash

IP=127.0.0.1
PORT=3045
echo
echo  $IP $PORT

echo
echo $(date +%Y-%m-%d,%H:%M:%s)

BODY="$(cat beis_req.txt)"

echo $BODY
echo
echo "--------START telnet---------"
echo

(sleep 1;echo "$BODY";echo;echo;sleep 300;)|telnet $IP $PORT
echo
echo "--------END-------------"

#
############################################################
