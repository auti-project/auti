#!/usr/bin/env bash

cd ../../clolc/benchmark || exit

go build -o clolc.out

LOG_DIR="../logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_init_epoch_bench.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

echo "CLOLC INIT EPOCH BENCH" >$LOG_FILE_DIR

for numOrg in 2 4 8 16 32 64 128 256; do
  ./clolc.out -phase i -numOrg $numOrg -numIter 10 | tee -a $LOG_FILE_DIR
  sleep 1
done

rm clolc.out
