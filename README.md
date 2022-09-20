[![CircleCI](https://circleci.com/gh/prometheus-community/smartctl_exporter.svg?style=svg)](https://circleci.com/gh/prometheus-community/smartctl_exporter)

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
      --smartctl.interval=60s  The interval between smarctl polls
      --smartctl.device=SMARTCTL.DEVICE ...  
                               The device to monitor (repeatable)
      --web.listen-address=":9633"  
                               Address to listen on for web interface and telemetry
      --web.telemetry-path="/metrics"  
                               Path under which to expose metrics
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
