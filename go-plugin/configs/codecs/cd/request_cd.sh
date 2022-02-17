#!/bin/bash
#5.2.2 CIMT000070-FB000.004
###########################################################
#
#

IP=127.0.0.1
PORT=2045
EXT_REF="$(date +%Y%m%d%H%M%S)$RANDOM"
echo
echo "=======5.2.2======"
echo  $IP $PORT
echo $EXT_REF ${#EXT_REF}
echo
echo $(date +%Y-%m-%d,%H:%M:%s)

BODY="$(sed 's/EXT_REQUEST_ID/'$EXT_REF'/g' cd_req.txt)"

BODY=$(printf "%010d" ${#BODY})${BODY}

echo $BODY
echo
echo "--------START telnet---------"
echo

#telnet $IP $PORT|(sleep 1;echo "$BODY";echo;echo "fff";sleep 30;) >> telnet_result.txt
(sleep 1;echo "$BODY";echo;echo;sleep 300;)|telnet $IP $PORT
echo
echo "--------END-------------"

#
############################################################
