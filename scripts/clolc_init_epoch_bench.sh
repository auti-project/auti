#!/usr/bin/env bash

cd ../clolc || exit

go build -o clolc.out

./clolc.out -phase i -numOrg 2 -numIter 10
./clolc.out -phase i -numOrg 4 -numIter 10
./clolc.out -phase i -numOrg 8 -numIter 10
./clolc.out -phase i -numOrg 16 -numIter 10
./clolc.out -phase i -numOrg 32 -numIter 10
./clolc.out -phase i -numOrg 64 -numIter 10
# ./clolc.out -phase i -numOrg 128 -numIter 10
# ./clolc.out -phase i -numOrg 256 -numIter 10
# ./clolc.out -phase i -numOrg 512 -numIter 10

rm clolc.out