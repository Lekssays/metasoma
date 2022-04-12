#!/bin/bash

echo "Welcome from $DISCOVERY_ADDRESS"

systemctl enable redis-server

service redis-server restart

export PATH=$PATH:/usr/local/go/bin

cd /declimnet/autopeering
rm *.pem autopeering
go mod tidy
go build
./autopeering server
