#! /usr/bin/env python3
# SPDX-License-Identifier: Apache-2.0
# Redact potentially sensitive information in smartctl JSON files
# This script does an in-place modification.
import json
import sys
import copy
import os

def main():
    for arg in sys.argv[1:]:
        print(arg)
        redact_one_file(arg)

def redact_one_file(filename):
    data = None
    tmpname = filename+".new"
    with open(filename, "r") as jsonFile:
        data = json.load(jsonFile)

    newdata = redact_data(data)

    with open(tmpname, "w") as jsonFile:
        json.dump(newdata, jsonFile, indent="\t", sort_keys=True)

    os.rename(tmpname, filename)

def mutate_nested_dict(d, keys, newvalue, if_present=False):
    # if_present=True: only mutate if the full key path exists.
    if len(keys) == 1:
        if not if_present or keys[0] in d:
            d[keys[0]] = newvalue
    else:
        k = keys[0]
        if k in d:
            mutate_nested_dict(d[k], keys[1:], newvalue, if_present=if_present)

REDACTED_STRING = 'REDACTED'
REDACTED_TIME_T = 1234567890
REDACTED_ASCTIME = "Fri Feb 13 23:31:30 2009 UTC" # TODO: generate from TIME_T, with UTC
REDACTED_HEX16_STR = '0x1234567890abcdef'
REDACTED_UINT32 = 1234567890

REDACT_FIELDS = [
    {'k': ['smartctl','platform_info'], 'v': REDACTED_STRING},
    {'k': ['smartctl','build_info'], 'v': REDACTED_STRING},
    {'k': ['serial_number'], 'v': REDACTED_STRING},
    {'k': ['firmware_version'], 'v': REDACTED_STRING},
    {'k': ['local_time', 'time_t'], 'v': REDACTED_TIME_T},
    {'k': ['local_time', 'asctime'], 'v': REDACTED_ASCTIME},
    {'k': ['logical_unit_id'], 'v': REDACTED_HEX16_STR},
    {'k': ['wwn','id'], 'v': REDACTED_UINT32},
    # TODO: how to redact /dev/sdX /dev/nvmeN ??
]

def redact_data(data):
    newdata = copy.deepcopy(data)
    for f in REDACT_FIELDS:
        #newval = str(f['v'])+str(f['k']) # for debugging
        newval = f['v']
        mutate_nested_dict(newdata, f['k'], newval, if_present=True)
    return newdata

if __name__ == '__main__':
    main()
