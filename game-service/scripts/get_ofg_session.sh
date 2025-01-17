#!/bin/bash

PLAYERID="gerard"
PASSWORD="tgTest987!"
CASINOID="c064fe0c135d6264"
GAMEID="OFG96"

TMPFILE="/tmp/xxx_mock_session_xxx.json"
PAYLOAD=$(jq -nc --arg p1 "$PLAYERID" --arg p2 "$PASSWORD" '{playerId:$p1,password:$p2}')
curl -q -X "POST" -H "Content-Type: application/json" "https://mcgw.dev.topgaming.team/mc/v1/login" -d "$PAYLOAD" > $TMPFILE 2>/dev/null
MOCKSESSION=$(jq -r .session $TMPFILE)
rm $TMPFILE
echo "session: $MOCKSESSION"

REDIRECT=$(curl "https://mcgw.dev.topgaming.team/mc/v1/launch?playerId=$PLAYERID&session=$MOCKSESSION&gameId=$GAMEID" 2>/dev/null)
TOKEN=${REDIRECT#*player_token=}
TOKEN=${TOKEN%&amp;*}
echo "token: $TOKEN"

PAYLOAD=$(jq -nc --arg p1 "$CASINOID" --arg p2 "$TOKEN" --arg p3 $PLAYERID --arg p4 $GAMEID '{casinoId:$p1,token:$p2,playerId:$p3,gameId:$p4}')
curl -q -X "POST" -H "Content-Type: application/json" "https://gw.dev.topgaming.team/ds/v1/auth" -d "$PAYLOAD" > $TMPFILE 2>/dev/null
SESSIONID=$(jq -r .sessionId $TMPFILE)
rm $TMPFILE
echo "sessionID: $SESSIONID"
