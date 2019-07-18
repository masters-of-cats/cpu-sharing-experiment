#!/bin/bash

set -e

APP_NAME=$1

function delete_cgroup {
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
}

trap cleanup EXIT

DIR=$(dirname $0)
APP_DIR="$DIR/../$1"

mkdir "$APP_DIR"
mkdir -p "/sys/fs/cgroup/cpu/good/$1"

echo $$ > "$APP_DIR/pidfile"
echo $$ > "/sys/fs/cgroup/cpu/good/$1/tasks"

while true; do
  if [ ! -f "$APP_DIR/spike" ]; then
    sleep 1
  else
    for i in $(seq 1000000); do
      true
    done
  fi
done
