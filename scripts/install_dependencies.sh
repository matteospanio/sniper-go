#!/usr/bin/env bash
# install_dependencies.sh
#
# author: Matteo Spanio <dev2@audioinnova.com>
# Program to install sniper-go dependencies

source scripts/utils.sh

function check_version {
    # $1: command to check version
    # $2: version to check
    OUT=$($1 --version)

    # if leading v is missing, add it
    if [[ $OUT =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        OUT="v$OUT"
    fi

    if [ $? -eq 0 ]; then
        ok_log "$1 is already installed"
        VERSION=$(echo $OUT | sed -E 's/^v([0-9]+)\.([0-9]+)\.([0-9]+)$/\1/')
        if [ $VERSION -lt $2 ]; then
            err_log "$1 version is less than $2, installing $1 $2"
            return 1
        else
            return 0
        fi
    else
        return 1
    fi
}

check_root

SNIPER_GO_DIR=$(pwd)
cd /tmp

# update repositories 
apt update

# install curl
apt install -y curl git-all make

# install git, node, npm

# check node version
check_version "node" 16
# if node version is less than 16, install node 16
if [ $? -eq 1 ]; then
    curl -fsSL https://deb.nodesource.com/setup_16.x | bash -
    apt install -y nodejs
    assert_error "Error installing node"
fi

# check npm version
check_version "npm" 6
# if npm version is less than 6, install npm
if [ $? -eq 1 ]; then
    apt install -y npm
    assert_error "Error installing npm"
fi

# install yarn
OUT=$(yarn --version)
if [ $? -eq 0 ]; then
    ok_log "Yarn is already installed"
else
    npm install -g yarn
    assert_error "Error installing yarn"
fi

# install sniper
OUT=$(sniper --help)
if [ $? -eq 0 ]; then
    ok_log "Sniper is already installed"
else
    git clone https://github.com/1N3/Sn1per
    cd Sn1per
    bash install.sh
    assert_error "Error installing sniper"
fi

# install go
OUT=$(go version)
if [ $? -eq 0 ]; then
    ok_log "Go is already installed"
else
    wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
fi

mkdir -p /usr/share/sniper/loot/workspace

cd $SNIPER_GO_DIR
