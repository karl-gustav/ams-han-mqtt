# Setup:

Setup the folders and service on remote machine:

    ./setup.sh <optional server name>

Deploy to server (run this locally):

    ./build.sh && ./deploy.sh <optional server name>

And check the logs that it went ok (have to be done on the remote):

    journalctl -u ams-han-mqtt # add -f to continuously print new entries
