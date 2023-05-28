#!/usr/bin/env bash

function clean_up() {
  echo "Clean up"
  ./fablo down
  rm -rf fablo-target
  docker volume prune -f
  docker network prune -f
  docker container prune -f
  rm lc_tx_id.log
  rm oc_tx_id.log
}