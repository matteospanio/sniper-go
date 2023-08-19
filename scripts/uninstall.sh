#!/usr/bin/env bash
# uninstall.sh
#
# author: Matteo Spanio <dev2@audioinnova.com>
# Uninstall sniper-go
source scripts/utils.sh
check_root

echo -e "$BOLDRED[>]$RESET Are you sure you want to uninstall sniper-go? (y/n)"
read -p ">>> " choice

if [ "$choice" != "y" ]; then
    err_log "Uninstall aborted"
    exit 1
fi

ok_log "Uninstalling sniper-go..."

rm -rf /usr/local/bin/sniper-go
rm -rf /usr/local/share/sniper-go

ok_log "Uninstall completed"