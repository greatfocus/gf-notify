#!/bin/sh
export PATH=$PATH:/usr/local/go/bin

# Building our app
GOOS=linux GOARCH=amd64 go build

# create service
sudo systemctl stop gf-notify 
sudo systemctl disable gf-notify  
sudo cp gf-notify.service /etc/systemd/system/gf-notify.service
systemctl daemon-reload

# start user 
sudo systemctl enable gf-notify 
sudo systemctl start gf-notify  