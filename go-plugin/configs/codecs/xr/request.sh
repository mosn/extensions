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
echo $EXTREF
echo
echo $(date +%Y-%m-%d,%H:%M:%s)

BODY='<Service><Header><ServiceCode>CIMT000070</ServiceCode><ChannelId>C48</ChannelId><ExternalReference>'$EXT_REF'</ExternalReference><OriginalChannelId>C49</OriginalChannelId><OriginalReference>06221113270051159201000092010000</OriginalReference><RequestTime>20210622111327543</RequestTime><Version>1.0</Version><RequestType>0</RequestType><Encrypt>0</Encrypt><TradeDate>20210617</TradeDate><RequestBranchCode>CN0010001</RequestBranchCode><RequestOperatorId>FB.ICP.X01</RequestOperatorId><RequestOperatorType>1</RequestOperatorType><TermType>00000</TermType><TermNo>0000000000</TermNo></Header><Body><Request><CustNo>3001504094</CustNo></Request></Body></Service>'

BODY=$(printf "%08d" ${#BODY})${BODY}

echo
echo "--------START invoke---------"
echo

echo "$BODY"

# (sleep 1;echo "$BODY";echo;echo;sleep 300;)|telnet $IP $PORT
(echo "$BODY"; sleep 15;) | nc -4 $IP $PORT
echo
echo "--------END-------------"

#
############################################################
