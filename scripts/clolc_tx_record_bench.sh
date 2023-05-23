#!/usr/bin/env bash

rm -rf fablo-target

# install fablo if not installed
[ -f ./fablo ] ||
curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o ./fablo && chmod +x ./fablo

export AUTI_ORG_GLOBAL_DIR=${PWD}

./fablo up fablo-config.yaml