#!/usr/bin/env bash

cd ../benchmark || exit

go build -o clolc.out

LOG_DIR="../logs"
if [ ! -d $LOG_DIR ]; then
    mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_tr_local_chain_submit.log"
if [ -f $LOG_FILE_DIR ]; then
    rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

# install fablo if not installed
[ -f ./fablo ] ||
curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o ./fablo && chmod +x ./fablo

export AUTI_LOCAL_CHAIN_DIR=${PWD}

function cleanup() {
    ./fablo down
    rm -rf fablo-target
    docker volume prune -f
    docker network prune -f
    docker container prune -f
}

# 1k test
cleanup
./fablo up fablo-local-chain-config.yaml
sleep 5
./clolc.out -phase tr -process local_submit -numTXs 1000 -numIter 10 | tee -a $LOG_FILE_DIR
sleep 1

for i in 10000 100000 1000000; do
    for j in {1..10}; do
        echo "No: $j" >>$LOG_FILE_DIR
        cleanup
        ./fablo up fablo-local-chain-config.yaml
        sleep 5
        ./clolc.out -phase tr -process local_submit -numTXs $i -numIter 1 | tee -a $LOG_FILE_DIR
        sleep 1
    done
done

rm clolc.out
