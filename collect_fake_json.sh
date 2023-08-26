#! /bin/bash

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# The original script used --xall but that doesn't work
# This matches the command in readSMARTctl()
smartctl_args="--json --info --health --attributes --tolerance=verypermissive --nocheck=standby --format=brief --log=error"

[[ ! -d "${script_dir}/debug" ]] && mkdir --parents "${script_dir}/debug"

if command -v jq >/dev/null; then
	devices=$(smartctl --scan --json | jq --raw-output '.devices[].name')
elif command -v yq >/dev/null; then
	devices=$(smartctl --scan --json | yq --unwrapScalar '.devices[].name')
elif command -v awk >/dev/null; then
	devices=$(smartctl --scan | awk '{ print($1) }')
else
	devices=$(smartctl --scan | cut -d ' ' -f 1)
fi

for device in $devices; do
	echo "Collecting data for '${device}'"
	# shellcheck disable=SC2086
	sudo smartctl $smartctl_args "${device}" > "${script_dir}/debug/$(basename "${device}").json"
done
