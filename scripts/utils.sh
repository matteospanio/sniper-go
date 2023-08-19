GREEN='\033[0;32m'
RED='\033[0;31m'
RESET='\e[0m'
BOLDRED='\033[1;31m'
BOLDGREEN='\033[1;32m'
BOLD='\033[1m'

function check_root() {
    if [[ $EUID -ne 0 ]]; then
        echo "$BOLD[+]$RESET This script must be run as root"
        exit 1
    fi
}

function ok_log() {
    echo -e "$BOLDGREEN[+]$RESET $1"
}

function err_log() {
    echo -e "$BOLDRED[-]$RESET $1"
}

