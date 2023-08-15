[![CircleCI](https://circleci.com/gh/prometheus-community/smartctl_exporter.svg?style=svg)](https://circleci.com/gh/prometheus-community/smartctl_exporter)
[![Container Repository on Quay](https://quay.io/repository/prometheuscommunity/smartctl-exporter/status "Container Repository on Quay")](https://quay.io/repository/prometheuscommunity/smartctl-exporter)

# smartctl_exporter
Export smartctl statistics to prometheus

Example output you can show in [EXAMPLE.md](EXAMPLE.md)

## Need more?
**If you need additional metrics - contact me :)**
**Create a feature request, describe the metric that you would like to have and attach exported from smartctl json file**

# Requirements
smartmontools >= 7.0, because export to json [released in 7.0](https://www.smartmontools.org/browser/tags/RELEASE_7_0/smartmontools/NEWS#L11)

# Configuration
## Command line options

The exporter will scan the system for available devices if no `--smartctl.device` flags are used.

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
                               Regexp of devices to exclude from automatic scanning. (mutually exclusive to
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
