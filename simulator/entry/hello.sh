#!/bin/bash

echo "Welcome from $CORE_PEER_IP_ADDRESS"

systemctl enable redis-server

service redis-server restart

cd /autopeering
go mod tidy
go build
./autopeering server