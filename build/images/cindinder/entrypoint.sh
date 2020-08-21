#!/bin/ash

/usr/local/bin/dockerd-entrypoint.sh &
while [ ! -S /var/run/docker.sock ]; do echo "waiting for docker..."; sleep 2; done
docker load -i /cinder/images/cinder.tar
