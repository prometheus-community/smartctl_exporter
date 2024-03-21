#! /bin/bash

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# Data directory to dump smartctl output
# This directory will be created if it doesn't exist
data_dir="${script_dir}/smartctl-data"

# The original script used --xall but that doesn't work
# This matches the command in readSMARTctl()
smartctl_args="--json --info --health --attributes --tolerance=verypermissive \
--nocheck=standby --format=brief --log=error"

# Determine the json query tool to use
if command -v jq >/dev/null; then
	json_tool="jq"
	json_args="--raw-output"
elif command -v yq >/dev/null; then
	json_tool="yq"
	json_args="--unwrapScalar"
else
	echo -e "One of 'yq' or 'jq' is required. Please try again after \
installing one of them"
	exit 1
fi

if [[ ! "${UID}" -eq 0 ]] && ! command -v sudo >/dev/null; then
	# Not root and sudo doesn't exist
	echo "sudo does not exist. Please run this as root"
	exit 1
fi

SUDO="sudo"
if [[ "${UID}" -eq 0 ]]; then
	# Don't use sudo if root
	SUDO=""
fi

[[ ! -d "${data_dir}" ]] && mkdir --parents "${data_dir}"

if [[ $# -ne 0 ]]; then
	devices="${1}"
else
	devices="$(smartctl --scan --json | "${json_tool}" "${json_args}" \
'.devices[].name')"
  mapfile -t devices <<< "${devices[@]}"
fi

for device in "${devices[@]}"
  do
	echo -n "Collecting data for '${device}'..."
	# shellcheck disable=SC2086
	data="$($SUDO smartctl ${smartctl_args} ${device})"
	type="$(echo "${data}" | "${json_tool}" "${json_args}" '.device.type')"
	family="$(echo "${data}" | "${json_tool}" "${json_args}" \
'select(.model_family != null) | .model_family | sub(" |/" ; "_" ; "g")')"
	model="$(echo "${data}" | "${json_tool}" "${json_args}" \
'.model_name | sub(" |/" ; "_" ; "g")')"
	device_name="$(basename "${device}")"
	echo -e "\tSaving to ${type}-${family:=null}-${model}-${device_name}.json"
	echo "${data}" > \
"${data_dir}/${type}-${family:=null}-${model}-${device_name}.json"
done
