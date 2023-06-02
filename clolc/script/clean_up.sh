#!/usr/bin/env bash

function clean_up() {
  docker ps -aq | xargs docker stop | xargs docker rm
  echo "Clean up"
  ./fablo prune
  rm -rf fablo-target
  docker volume prune -f
  docker network prune -f
  docker container prune -f
  rm -rf wallet
  rm -rf keystore
  rm lc_tx_id.log
  rm oc_tx_id.log
}

clean_up