#!/usr/bin/env bash

CURR_DIR=${PWD}

cd $CURR_DIR/1_initialization || exit
./bench.sh
sleep 5

cd $CURR_DIR/2_transaction_record || exit
# ./1_bench_local_chain_submit.sh
# sleep 5
# ./2_bench_local_chain_read.sh
# sleep 5
# ./3_bench_org_chain_submit.sh
# sleep 5
# ./4_bench_org_chain_read.sh
# sleep 5
./5_bench_commit.sh
sleep 5
./6_bench_accumulate.sh
sleep 5

cd $CURR_DIR/3_consistency_examination || exit
./1_bench_acc_commit.sh
sleep 5
./2_bench_cal_b.sh
sleep 5
./3_bench_cal_c.sh
sleep 5
./4_bench_cal_d.sh
sleep 5
./5_bench_encrypt.sh
sleep 5
# ./6_bench_aud_chain_submit.sh
# sleep 5
# ./7_bench_aud_chain_read.sh
# sleep 5
./8_bench_check.sh