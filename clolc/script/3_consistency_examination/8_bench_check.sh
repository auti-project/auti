#!/usr/bin/env bash

source ../clean_up.sh

HOME_DIR="../.."
cd $HOME_DIR || exit

go build -o clolc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_ce_check.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

for j in {1..10}; do
  echo "No: $j" >>$LOG_FILE_DIR
  ./clolc.out -phase ce -process check -numIter 1 | tee -a $LOG_FILE_DIR
  sleep 1
done

rm clolc.out