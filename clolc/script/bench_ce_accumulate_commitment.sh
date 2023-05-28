#!/usr/bin/env bash

source ./clean_up.sh

cd ../benchmark || exit

clean_up

go build -o clolc.out

LOG_DIR="../logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_ce_accumulate_commitment.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

#for i in 1000 10000 100000 1000000; do
for j in {1..10}; do
  echo "No: $j" >>$LOG_FILE_DIR
  ./clolc.out -phase ce -process acc_commit -numOrg 2 -numIter 1 | tee -a $LOG_FILE_DIR
  sleep 1
done
#done

rm clolc.out
