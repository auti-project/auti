#!/usr/bin/env bash


# install fablo if not installed
[ -f ./fablo ] ||
curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o ./fablo && chmod +x ./fablo

export AUTI_LOCAL_CHAIN_DIR=${PWD}

./fablo down
rm -rf fablo-target
./fablo up fablo-config.yaml

cd ../benchmark || exit
go build -o tx_record.out
./tx_record.out