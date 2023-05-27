#!/usr/bin/env bash

cd ../benchmark || exit

./fablo down
rm -rf fablo-target
docker volume prune -f
docker network prune -f
docker container prune -f
rm lc_tx_id.log