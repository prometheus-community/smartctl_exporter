#!/bin/sh

DISKS="$(lsblk -l|egrep -oe '(sd[a-z]|nvme[0-9])'|sed -e 's/^/  - \/dev\//'| uniq)"
echo "$DISKS" >> smartctl_exporter.yaml

/bin/smartctl_exporter -config=/smartctl_exporter.yaml
