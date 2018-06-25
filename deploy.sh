#! /bin/bash
set -e

SERVER=${1:-smart}

CMD="tar xzf - -C /srv/ams-han-mqtt/\
 && echo Replacing executable \
 && test -f /srv/ams-han-mqtt/ams-han-mqtt \
    && mv /srv/ams-han-mqtt/ams-han-mqtt{,.old} || true \
 && mv /srv/ams-han-mqtt/ams-han-mqtt{.new,} \
 && echo Restarting service \
 && systemctl daemon-reload \
 && service ams-han-mqtt restart \
 && echo Checking status \
 && service ams-han-mqtt status \
"

echo 'Running command on "'${SERVER}'":' $CMD
tar czf - ams-han-mqtt.new ams-han-mqtt.service | ssh $SERVER $CMD
