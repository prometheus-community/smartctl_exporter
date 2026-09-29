# smartctl-mixin

A Prometheus / Grafana [mixin](https://github.com/monitoring-mixins/docs) for the
[`smartctl_exporter`](https://github.com/prometheus-community/smartctl_exporter).

It packages reusable, configurable alerting rules and a Grafana dashboard
describing SMART disk health across all scraped instances.

**Simply download them** from the [pre-compiled output directory](./output/) or modify your own [config file](./config.libsonnet) and follow the build instructions.

## Build

Requires `jsonnet` and `jsonnetfmt`; `jb` for dependency management and
`promtool` (optional) for rule validation are recommended:

```
go install github.com/google/go-jsonnet/cmd/jsonnet@latest
go install github.com/google/go-jsonnet/cmd/jsonnetfmt@latest
go install github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb@latest
```

Ensure `$(go env GOPATH)/bin` is on your `$PATH` (`~/go/bin` by default).

This mixin currently has no external Jsonnet dependencies, so `jb install` is
optional. If you add a `jsonnetfile.json` dependency (e.g. grafonnet), run:

```
jb install
```

Generate the Prometheus rule files and the Grafana dashboard:

```
make          # fmt + alerts + dashboards + lint
# or individually:
make output/smartctl_alerts.yaml
make output/dashboards
```

`output/` will contain `smartctl_alerts.yaml` and `smartctl-overview.json`,
ready to be loaded.

## Alerting rules

| alert | severity | trigger |
| --- | --- | --- |
| `SmartDeviceTemperatureWarning` | warning | avg temp > 60 °C (configurable) |
| `SmartDeviceTemperatureCritical` | critical | max temp > 70 °C (configurable) |
| `SmartDeviceTemperatureOverTripValue` | critical | temp exceeds drive trip value |
| `SmartDeviceTemperatureNearingTripValue` | warning | temp > 80 % of trip value |
| `SmartStatus` | critical | `smartctl_device_smart_status != 1` probe is reporting issues |
| `SmartCriticalWarning` | critical | `smartctl_device_critical_warning > 0` |
| `SmartMediaErrors` | critical | `smartctl_device_media_errors > 0` |
| `SmartWearoutIndicator` | critical | available spare below threshold |
| `SmartDeviceSSDWearingOut` | warning | SSD Media Wearout Indicator < 10 |
| `SmartErrorLogs` | critical | `smartctl_device_error_log_count > 0` |

