#! /bin/bash

tar czf - ams_han_mqtt | \
ssh smart 'tar xzf - -C .'
