#!/usr/bin/env bash

source ./clean_up.sh

cd ../benchmark || exit

clean_up

go build -o clolc.out

LOG_DIR="../logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_ce_compute_c.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

for j in {1..10}; do
  echo "No: $j" >>$LOG_FILE_DIR
  ./clolc.out -phase ce -process cal_c -numOrg 2 -numIter 1 | tee -a $LOG_FILE_DIR
  sleep 1
done

rm clolc.out