#!/bin/bash
#5.2.2 CIMT000070-FB000.004
###########################################################
#
#

IP=127.0.0.1
PORT=2045
echo
echo  $IP $PORT
echo
echo $(date +%Y-%m-%d,%H:%M:%s)

BODY="$(cat cd_req.txt)"

BODY=$(printf "%010d" ${#BODY})${BODY}

echo $BODY
echo "--------START invoke---------"

#telnet $IP $PORT|(sleep 1;echo "$BODY";echo;echo "fff";sleep 30;) >> telnet_result.txt
(sleep 1;echo "$BODY";) | nc $IP $PORT
echo
echo "--------END-------------"

#
############################################################
