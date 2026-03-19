#! /bin/bash
#
# Modified 2026-03-19: Added MegaRAID RAID controller support
# Changes:
#   - Removed /dev/bus device filtering to allow RAID controller devices
#   - Added device type extraction from smartctl scan to pass correct -d flag
#   - Added support for scsi_model_name in addition to model_name for SCSI devices

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# Data directory to dump smartctl output
# This directory will be created if it doesn't exist
data_dir="${script_dir}/smartctl-data"

# The original script used --xall but that doesn't work
# This matches the command in readSMARTctl()
smartctl_args="--json --info --health --attributes --tolerance=verypermissive \
--nocheck=standby --format=brief --log=error"

# Ignore these devices (empty by default to include all devices including /dev/bus for RAID controllers)
# Example: smartctl_ignore_dev_regex="^(/dev/sg)"
smartctl_ignore_dev_regex=""

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
	mapfile -t devices <<< "${devices[@]}"
else
	# Get both device name and type from smartctl scan
	scan_output="$(smartctl --scan --json)"
	devices=()
	while IFS= read -r device_json; do
		dev_name=$(echo "${device_json}" | "${json_tool}" "${json_args}" '.name')
		dev_type=$(echo "${device_json}" | "${json_tool}" "${json_args}" '.type')
		dev_label="${dev_name}"
		# Filter out ignored devices
		if [[ -n "${smartctl_ignore_dev_regex}" ]] && [[ "${dev_label}" =~ ${smartctl_ignore_dev_regex} ]]; then
			continue
		fi
		devices+=("${dev_name};${dev_type}")
	done < <(echo "${scan_output}" | "${json_tool}" "${json_args}" -c '.devices[]')
fi

for device_spec in "${devices[@]}"
  do
	# Split device spec into name and type (format: /dev/sdX;type or just /dev/sdX)
	IFS=';' read -r device device_type <<< "${device_spec}"
	[[ -z "${device_type}" ]] && device_type="auto"

	echo -n "Collecting data for '${device}' with type '${device_type}'..."
	# shellcheck disable=SC2086
	data="$($SUDO smartctl ${smartctl_args} --device=${device_type} ${device})"
	# Accommodate a smartmontools pre-7.3 bug
	data=${data#"  Pending defect count:"}
	type="$(echo "${data}" | "${json_tool}" "${json_args}" '.device.type')"
	family="$(echo "${data}" | "${json_tool}" "${json_args}" \
'select(.model_family != null) | .model_family | sub(" |/" ; "_" ; "g")
 | sub("\"|\\(|\\)" ; "" ; "g")')"
	# SCSI devices use scsi_model_name instead of model_name
	model="$(echo "${data}" | "${json_tool}" "${json_args}" \
'(.model_name // .scsi_model_name // "unknown") | sub(" |/" ; "_" ; "g") | sub("\"|\\(|\\)" ; "" ; "g")')"
	device_name="$(basename "${device}")"
	# For megaraid devices, basename returns "0", so use a more descriptive name
	if [[ "${device_type}" =~ ^(megaraid|sat\+megaraid) ]]; then
		device_name="${device_name}_${device_type//,/_}"
	fi
	echo -e "\tSaving to ${type}-${family:=null}-${model}-${device_name}.json"
	echo "${data}" > \
"${data_dir}/${type}-${family:=null}-${model}-${device_name}.json"
done
