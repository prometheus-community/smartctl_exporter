#! /bin/bash

for device in $(smartctl --scan | awk '{ print $1}')
do
  smartctl --json --xall $device | jq > $(basename $device).json
done
