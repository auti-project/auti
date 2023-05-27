#!/usr/bin/env bash

cd ../benchmark || exit

go build -o clolc.out

LOG_DIR="../logs"
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

function cleanup() {
    ./fablo down
    rm -rf fablo-target
    docker volume prune -f
    docker network prune -f
    docker container prune -f
    rm oc_tx_id.log
}

for i in 1000 10000 100000 1000000; do
    cleanup
    ./fablo up fablo-org-chain-config.yaml
    sleep 5
    ./clolc.out -phase tr -process org_prepare -numTXs $i | tee -a $LOG_FILE_DIR
    for j in {1..10}; do
        echo "No: $j" >>$LOG_FILE_DIR
        ./clolc.out -phase tr -process org_read -numTXs $i -numIter 1 | tee -a $LOG_FILE_DIR
        sleep 1
    done
done

for i in 1000 10000 100000 1000000; do
    cleanup
    ./fablo up fablo-org-chain-config.yaml
    sleep 5
    ./clolc.out -phase tr -process org_prepare -numTXs $i | tee -a $LOG_FILE_DIR
    for j in {1..10}; do
        echo "No: $j" >>$LOG_FILE_DIR
        ./clolc.out -phase tr -process org_read_all -numTXs $i -numIter 1 | tee -a $LOG_FILE_DIR
        sleep 1
    done
done

rm clolc.out
