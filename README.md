[![CircleCI](https://circleci.com/gh/Sheridan/smartctl_exporter.svg?style=svg)](https://circleci.com/gh/Sheridan/smartctl_exporter)

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
* `--config=/path/to/file.yaml`: Path to configuration file, defaulr `/etc/smartctl_exporter.yaml`
* `--verbose`: verbosed log, default no
* `--debug`: Debug logging, default no
* `--version`: Show version and exit

## Configuration file
Example content:
```
smartctl_exporter:
  bind_to: "[::1]:9633"
  url_path: "/metrics"
  fake_json: no
  smartctl_location: /usr/sbin/smartctl
  collect_not_more_than_period: 120s
  devices:
  - /dev/sda
  - /dev/sdb
  - /dev/sdc
  - /dev/sdd
  - /dev/sde
  - /dev/sdf
```
`fake_json` used for debugging.
