#! /bin/bash

for device in $(smartctl --scan | awk '{ print $1}')
do
  smartctl --json --xall $device | jq > debug/$(basename $device).json
done
