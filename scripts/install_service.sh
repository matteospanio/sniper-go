#!/usr/bin/env bash
# install_service.sh
#
# author: Matteo Spanio <dev2@audioinnova.com>
# Program to install sniper-go as a systemd service

SERVICE="sniper-go.service"

cp $SERVICE /etc/systemd/system/$SERVICE
systemctl daemon-reload
systemctl enable $SERVICE
echo "Sniper-go service installed"

echo "Do you want to start the service now? [y/n]"
read answer
if [ "$answer" == "y" ]; then
    systemctl start $SERVICE
    echo "Sniper-go service started"
fi
