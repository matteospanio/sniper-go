#!/usr/bin/env bash
# install_dependencies.sh
#
# author: Matteo Spanio <dev2@audioinnova.com>
# Program to install sniper-go dependencies

SNIPER_GO_DIR=$(pwd)

# update repositories 
apt update

# install git, node, npm
curl -fsSL https://deb.nodesource.com/setup_16.x | bash -
apt install -y nodejs npm git-all make

# install yarn
if [ ! -f /usr/bin/yarn ]; then
    npm install -g yarn
fi

# install go
if [ ! -f /usr/local/go/bin/go ]; then
    wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
fi

# install sniper
OUT=$(sniper --help)
if [ $? -eq 0 ]; then
    echo "Sniper is already installed"
else
    cd /tmp
    git clone https://github.com/1N3/Sn1per
    cd Sn1per
    bash install.sh
fi

cd $SNIPER_GO_DIR
