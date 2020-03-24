#! /bin/bash
set -e

SERVER=${1:-smart}

CMD="tar xzf - -C /tmp/ \
 && echo Making directory \
 && sudo mkdir -p /srv/ams-han-mqtt/ \
 && sudo chown \$USER /srv/ams-han-mqtt/ \
 && echo Copying service \
 && sudo cp /tmp/ams-han-mqtt.service /srv/ams-han-mqtt/ \
 && echo Enabling service \
 && sudo systemctl enable /srv/ams-han-mqtt/ams-han-mqtt.service \
"

echo 'Running command on "'${SERVER}'":' $CMD
tar czf - ams-han-mqtt.service |ssh $SERVER $CMD
