[![CircleCI](https://circleci.com/gh/prometheus-community/smartctl_exporter.svg?style=svg)](https://circleci.com/gh/prometheus-community/smartctl_exporter)
[![Container Repository on Quay](https://quay.io/repository/prometheuscommunity/smartctl-exporter/status "Container Repository on Quay")](https://quay.io/repository/prometheuscommunity/smartctl-exporter)

# smartctl_exporter
Export smartctl statistics to prometheus

Example output you can show in [EXAMPLE.md](EXAMPLE.md)

## Need more?
**If you need additional metrics - contact me :)**

**Create a feature request, describe the metric that you would like to have and
attach exported from smartctl json file**

# Requirements
`smartmontools` >= 7.0, because export to json [released in 7.0](https://www.smartmontools.org/browser/tags/RELEASE_7_0/smartmontools/NEWS#L11)

# Configuration
## Command line options

The exporter will scan the system for available devices if no `--smartctl.device`
flags are used.

```
usage: smartctl_exporter [<flags>]

Flags:
  -h, --help                   Show context-sensitive help (also try --help-long and --help-man).
      --smartctl.path="/usr/sbin/smartctl"  
                               The path to the smartctl binary
      --smartctl.interval=60s  The interval between smartctl polls
      --smartctl.rescan=10m    The interval between rescanning for new/disappeared devices. If the interval is smaller than 1s no
                               rescanning takes place. If any devices are configured with smartctl.device also no rescanning takes
                               place.
      --smartctl.device=SMARTCTL.DEVICE ...  
                               The device to monitor (repeatable)
      --smartctl.device-exclude=""
                               Regexp of devices to exclude from automatic scanning. (mutually exclusive to
                               device-include)
      --smartctl.device-include=""
                               Regexp of devices to include in automatic scanning. (mutually exclusive to
                               device-exclude)
      --web.telemetry-path="/metrics"  
                               Path under which to expose metrics
      --web.systemd-socket     Use systemd socket activation listeners instead of port listeners (Linux only).
      --web.listen-address=:9633 ...
                               Addresses on which to expose metrics and web interface. Repeatable for multiple
                               addresses.
      --web.config.file=""     [EXPERIMENTAL] Path to configuration file that can enable TLS or authentication.
      --log.level=info         Only log messages with the given severity or above. One of: [debug, info, warn,
                               error]
      --log.format=logfmt      Output format of log messages. One of: [logfmt, json]
      --version                Show application version.
```

## TLS and basic authentication

This exporter supports TLS and basic authentication.

To use TLS and/or basic authentication, you need to pass a configuration file
using the `--web.config.file` parameter. The format of the file is described
[in the exporter-toolkit repository](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md).

## Example of running in Docker

Minimal functional `docker-compose.yml`:
```yaml
version: "3"

services:
  smartctl-exporter:
    image: prometheuscommunity/smartctl-exporter
    privileged: true
    user: root
    ports:
      - "9633:9633"
```

# Troubleshooting
## Troubleshooting data inconsistencies
`smartmon_exporter` uses the JSON output from `smartctl` to provide the data to
Prometheus. If the data is incorrect, look at the data from `smartctl` to
determine if the issue should be reported upstream to smartmontools or to this
repo. In general, the `smartctl_exporter` should not modify the data in flight.
If the data is missing from `smartctl`, it should not be in `smartctl_exporter`.
If the data from `smartctl` is incorrect, it should be reported upstream.
Requests for `smartctl_exporter` to "fix" incorrect data where `smartctl` is
reporting incorrect data will be closed. The grey area is when invalid or
missing data from smartctl is causing multiple invalid or incorrect data
in `smartctl_exporter`. This could happen if the data is used in a calculation
for other data. This will need to be researched on a case by case basis.

| -                         | smartctl valid              | smartctl missing                                | smartctl invalid/incorrect       |
|---------------------------|-----------------------------|-------------------------------------------------|----------------------------------|
| smartctl_exporter valid   | all good                    | N/A                                             | N/A                              |
| smartctl_exporter missing | issue for smartctl_exporter | report upstream to smartmontools                | report upstream to smartmontools |
| smartctl_exporter invalid | issue for smartctl_exporter | issue for smartctl_exporter and report upstream | report upstream to smartmontools |

### smartctl output vs smartctl_exporter output

The S.M.A.R.T. attributes are mapped in
[smartctl.go](https://github.com/prometheus-community/smartctl_exporter/blob/master/smartctl.go).
Each function has a `prometheus.MustNewConstMetric` or similar function with the
first parameter being the metric name. Find the metric name in
[metrics.go](https://github.com/prometheus-community/smartctl_exporter/blob/master/metrics.go)
to see how the exporter displays the information. This may sound technical, but
it's crucial for understanding how data flows from `smartctl` to
`smartctl_exporter` to Prometheus.

If the data looks incorrect, check the
[Smartmontools Frequently Asked Questions (FAQ)](https://www.smartmontools.org/wiki/FAQ).
It's likely your question may already have an answer. If you still have
questions, open an [issue]().

## Gathering smartctl data
Follow these steps to gather smartctl data for troubleshooting purposes. If you
have unique drives/data/edge cases and would like to "donate" the data, open a
PR with the redacted JSON files.

1. Run `scripts/collect-smartctl-json.sh` to export all drives to a
   `smartctl-data` directory (created in the current directory).
2. Run `scripts/redact_fake_json.py` to redact sensitive data.
3. Provide the JSON file for the drive in question.

```bash
cd scripts
./collect-smartctl-json.sh
./redact-fake-json.py smartctl-data/*.json
```

## Run smartctl_exporter using JSON data
The `smartctl_exporter` can be run using local JSON data. The device names are
pulled from actual devices in the machine while the data is redirected to the
`debug` directory. Save the JSON data in the `debug` directory using the actual
device names using a 1:1 ratio. If you have 3 devices, `sda`, `sdb` and `sdc`,
the `smartctl_exporter` will expect 3 files: `debug/sda.json`, `debug/sdb.json`
and `debug/sdc.json`.

Once the "fake devices" (JSON files) are in place, run the exporter passing the
hidden `--smartctl.fake-data` switch on the command line. The port is specified
to prevent conflicts with an existing `smartctl_exporter` on the default port.

```bash
smartctl_exporter --web.listen-address 127.0.0.1:19633 --smartctl.fake-data
```

# FAQ
## How do I run `smartctl_exporter` against a JSON file?

If you're helping someone else, request the output of the `smartctl` command
above. Feed this into the `smartctl_exporter` using the
hidden `--smartctl.fake-data` flag. If a `smartctl_exporter` is already running,
use a different port; in this case, it's `19633`. Run `collect_fake_json.sh`
first to collect the JSON files for **your** devices. Copy the requested JSON
file into one of the fake files. After starting the exporter, you can query it
to see the data generated.

```bash
# Dump the JSON files for your devices into debug/
./collect_fake_json.sh

# copy the test JSON into one of the files in debug/
cp extracted-from-above-sda.json debug/sda.json

# Make sure you have the latest version
go build
# Use a different port in case smartctl_exporter is already running
sudo ./smartctl_exporter --web.listen-address=127.0.0.1:19633 --log.level=debug --smartctl.fake-data

# Use curl with grep
curl --silent 127.0.0.1:19633/metrics | grep -i nvme
# Or xh with ripgrep
xh --body :19633/metrics | rg nvme
```

## Why is root required? Can't I add a user to the "disk" group?

A blogger had the same question and opened a ticket on smartmontools. This is
their response. `smartctl` needs to be run as root.

[RFE: add O_RDRW mode for sat/scsi/ata devices](https://www.smartmontools.org/ticket/1064)

> According to function `blk_verify_command()` from current kernel sources
> (see [​block/scsi_ioctl.c](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/block/scsi_ioctl.c)),
> O_RDONLY or O_RDWR make no difference if device was opened as root (or with
> CAP_SYS_RAWIO).
>
> The SCSI commands listed in function `blk_set_cmd_filter_defaults()` show
> that some of the `smartctl -d scsi` functionality might work with O_RDONLY
> for non-root users. Some more might work with O_RDWR.
>
> But `smartctl -d sat` (to access SATA devices) won't work at all because the
> SCSI commands ATA_12 and ATA_16
> (see [​scsi_proto.h](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/scsi/scsi_proto.h))
> are **always blocked for non-root users**.

## What about my NVMe drive?
From the smartmontools FAQ: [My NVMe drive is not in the smartctl/smartd database](https://www.smartmontools.org/wiki/FAQ#MyNVMedriveisnotinthesmartctlsmartddatabase)
> SCSI/SAS and NVMe drives do not provide ATA/SATA-like SMART Attributes.
> Therefore the drive database does not contain any entries for these drives.
> This may change in the future as some drives provide similar info via vendor
> specific commands (see ticket #870).

smartmontools also has a [wiki page for NVMe](https://www.smartmontools.org/wiki/NVMe_Support) devices.

## How do I report upstream to smartmontools?
Check their FAQ: [How to create a bug report](https://www.smartmontools.org/wiki/FAQ#Howtocreateabugreport).
