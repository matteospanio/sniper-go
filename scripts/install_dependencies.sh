#!/usr/bin/env bash
# install_dependencies.sh
#
# author: Matteo Spanio <dev2@audioinnova.com>
# Program to install sniper-go dependencies

source scripts/utils.sh

check_root

SNIPER_GO_DIR=$(pwd)
cd /tmp

# update repositories 
apt update

# install curl
apt install -y curl git-all make

# install git, node, npm

OUT=$(node --version)
# get node version major
NODE_VERSION=$(echo $OUT | sed -E 's/^v([0-9]+)\.([0-9]+)\.([0-9]+)$/\1/')

# if node version is less than 16, install node 16
if [ $NODE_VERSION -lt 16 ]; then
    err_log "Node version is less than 16, installing node 16"
    curl -fsSL https://deb.nodesource.com/setup_16.x | bash -
    apt install -y nodejs npm
fi

# install yarn
OUT=$(yarn --version)
if [ $? -eq 0 ]; then
    ok_log "Yarn is already installed"
else
    npm install -g yarn
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

# install sniper
OUT=$(sniper --help)
if [ $? -eq 0 ]; then
    ok_log "Sniper is already installed"
else
    git clone https://github.com/1N3/Sn1per
    cd Sn1per
    bash install.sh
fi

cd $SNIPER_GO_DIR
