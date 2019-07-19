#!/bin/bash

set -e

APP_NAME=$1

function delete_cgroup {
  if [ ! -d "/sys/fs/cgroup/cpu/$1/$APP_NAME" ]; then
    return
  fi

  while true; do
   if rmdir "/sys/fs/cgroup/cpu/$1/$APP_NAME" 2> /dev/null; then
     break
   fi
   sleep .1
 done
}

function cleanup {
 rm -rf "$APP_DIR"
 echo $$ > /sys/fs/cgroup/cpu/tasks
 delete_cgroup "good"
 delete_cgroup "bad"
 exit
}

function eat_cpu {
  while true; do
    if [ ! -f "$1/spike" ]; then
      sleep 5
    else
      for i in $(seq 1000000); do
        true
      done
    fi
  done
}

export -f eat_cpu

trap cleanup EXIT

DIR=$(dirname $0)
APP_DIR="$DIR/../$1"

mkdir "$APP_DIR"
mkdir -p "/sys/fs/cgroup/cpu/good/$1"

echo $$ > "$APP_DIR/pidfile"
echo $$ > "/sys/fs/cgroup/cpu/good/$1/tasks"

parallel -n0 eat_cpu ::: $APP_NAME $APP_NAME
