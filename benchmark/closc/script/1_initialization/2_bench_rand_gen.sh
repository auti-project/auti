#!/usr/bin/env bash

HOME_DIR="../.."
cd $HOME_DIR || exit

source ./script/clean_up.sh

go build -o closc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/closc_in_rand_gen.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

for numRand in 1000 10000 100000 1000000; do
  for j in {1..10}; do
    echo "No: $j" >>$LOG_FILE_DIR
    ./closc.out -phase in -process rand_gen -num $numRand -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 1
  done
done

rm clolc.out