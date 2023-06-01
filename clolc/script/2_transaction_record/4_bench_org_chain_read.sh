#!/usr/bin/env bash

source ../clean_up.sh

HOME_DIR="../.."
cd $HOME_DIR || exit

go build -o clolc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_tr_org_chain_read.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

# install fablo if not installed
[ -f ./fablo ] ||
  curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o ./fablo && chmod +x ./fablo

export AUTI_ORG_CHAIN_DIR=${PWD}

FABLO_ORG_CHAIN_CONFIG="./config/fablo-org-chain-config.yaml"
clean_up
./fablo up $FABLO_ORG_CHAIN_CONFIG
sleep 5
TOTLA_TXS=0
for i in 1000 9000 90000 900000; do
  ./clolc.out -phase tr -process org_prepare -numTXs $i | tee -a $LOG_FILE_DIR
  TOTLA_TXS=$((TOTLA_TXS + i))
  sleep 5
  for j in {1..11}; do
    echo "No: $j" >>$LOG_FILE_DIR
    ./clolc.out -phase tr -process org_read -numTXs $TOTLA_TXS -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 5
    ./clolc.out -phase tr -process org_read_all -numTXs $TOTLA_TXS -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 5
  done
done

rm clolc.out
