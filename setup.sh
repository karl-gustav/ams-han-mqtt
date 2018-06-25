#! /bin/bash
set -e

SERVER=${1:-smart}

CMD="tar xzf - -C /tmp/ \
 && echo Making directory \
 && mkdir -p /srv/ams-han-mqtt/ \
 && echo Copying service \
 && cp /tmp/ams-han-mqtt.service /srv/ams-han-mqtt/ \
 && echo Enabling service \
 && systemctl enable /srv/ams-han-mqtt/ams-han-mqtt.service \
"

echo 'Running command on "'${SERVER}'":' $CMD
tar czf - ams-han-mqtt.service |ssh $SERVER $CMD
