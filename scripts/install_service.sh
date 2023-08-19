#!/usr/bin/env bash
# install_service.sh
#
# author: Matteo Spanio <dev2@audioinnova.com>
# Program to install sniper-go as a systemd service

source scripts/utils.sh
check_root

SERVICE="sniper-go.service"

cp $SERVICE /etc/systemd/system/$SERVICE
systemctl daemon-reload
systemctl enable $SERVICE
ok_log "Sniper-go service installed"

echo "$BOLD[>]$RESET Do you want to start the service now? [y/n]"
read answer
if [ "$answer" == "y" ]; then
    systemctl start $SERVICE
    if [ $? -ne 0 ]; then
        err_log "Error starting sniper-go service"
    else
    ok_log "Sniper-go service started"
fi
