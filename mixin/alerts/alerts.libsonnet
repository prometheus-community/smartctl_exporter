{
  local c = $._config,
  local sel = c.smartctlExporterSelector,

  prometheusAlerts+:: {
    groups+: [
      {
        name: 'smartctl',
        interval: '%s' % c.ruleEvaluationInterval,
        rules: [
          // ----------------------------------------------------------------
          // Temperature
          // ----------------------------------------------------------------
          {
            alert: 'SmartDeviceTemperatureWarning',
            expr: |||
              (
                avg_over_time(smartctl_device_temperature{%(sel)s, temperature_type="current"}[%(temperatureAvgWindow)s])
                  unless on (instance, device) smartctl_device_temperature{%(sel)s, temperature_type="drive_trip"}
              ) > %(temperatureWarningThreshold)d
            ||| % (c { sel: sel }),
            'for': '%s' % c.alertFor,
            labels: { severity: 'warning' },
            annotations: {
              summary: 'SMART device temperature warning (instance {{ $labels.instance }})',
              description: 'Device temperature warning on {{ $labels.instance }} drive {{ $labels.device }} over %(temperatureWarningThreshold)d°C. VALUE = {{ $value }}. LABELS = {{ $labels }}' % c,
            },
          },
          {
            alert: 'SmartDeviceTemperatureCritical',
            expr: |||
              (
                max_over_time(smartctl_device_temperature{%(sel)s, temperature_type="current"}[%(temperatureAvgWindow)s])
                  unless on (instance, device) smartctl_device_temperature{%(sel)s, temperature_type="drive_trip"}
              ) > %(temperatureCriticalThreshold)d
            ||| % (c { sel: sel }),
            'for': '%s' % c.alertFor,
            labels: { severity: 'critical' },
            annotations: {
              summary: 'SMART device temperature critical (instance {{ $labels.instance }})',
              description: 'Device temperature critical on {{ $labels.instance }} drive {{ $labels.device }} over %(temperatureCriticalThreshold)d°C. VALUE = {{ $value }}. LABELS = {{ $labels }}' % c,
            },
          },
          {
            alert: 'SmartDeviceTemperatureOverTripValue',
            expr: |||
              max_over_time(smartctl_device_temperature{%(sel)s, temperature_type="current"}[%(temperatureMaxWindow)s])
                > on (device, instance) smartctl_device_temperature{%(sel)s, temperature_type="drive_trip"}
            ||| % (c { sel: sel }),
            'for': '%s' % c.alertFor,
            labels: { severity: 'critical' },
            annotations: {
              summary: 'SMART device temperature over trip value (instance {{ $labels.instance }})',
              description: 'Device temperature over trip value on {{ $labels.instance }} drive {{ $labels.device }}. VALUE = {{ $value }}. LABELS = {{ $labels }}',
            },
          },
          {
            alert: 'SmartDeviceTemperatureNearingTripValue',
            expr: |||
              max_over_time(smartctl_device_temperature{%(sel)s, temperature_type="current"}[%(temperatureMaxWindow)s])
                > on (device, instance) (smartctl_device_temperature{%(sel)s, temperature_type="drive_trip"} * %(temperatureNearingTripRatio)f)
            ||| % (c { sel: sel }),
            'for': '%s' % c.alertFor,
            labels: { severity: 'warning' },
            annotations: {
              summary: 'SMART device temperature nearing trip value (instance {{ $labels.instance }})',
              description: 'Device temperature nearing trip value on {{ $labels.instance }} drive {{ $labels.device }}. VALUE = {{ $value }}. LABELS = {{ $labels }}',
            },
          },

          // ----------------------------------------------------------------
          // SMART status / critical warning / media errors
          // ----------------------------------------------------------------
          {
            alert: 'SmartStatus',
            expr: 'smartctl_device_smart_status{%(sel)s} != 1' % { sel: sel },
            'for': '%s' % c.alertFor,
            labels: { severity: 'critical' },
            annotations: {
              summary: 'SMART status (instance {{ $labels.instance }})',
              description: 'Device has a SMART status failure on {{ $labels.instance }} drive {{ $labels.device }}. VALUE = {{ $value }}. LABELS = {{ $labels }}',
            },
          },
          {
            alert: 'SmartCriticalWarning',
            expr: 'smartctl_device_critical_warning{%(sel)s} > 0' % { sel: sel },
            'for': '%s' % c.alertFor,
            labels: { severity: 'critical' },
            annotations: {
              summary: 'SMART critical warning (instance {{ $labels.instance }})',
              description: 'Disk controller has critical warning on {{ $labels.instance }} drive {{ $labels.device }}. VALUE = {{ $value }}. LABELS = {{ $labels }}',
            },
          },
          {
            alert: 'SmartMediaErrors',
            expr: 'smartctl_device_media_errors{%(sel)s} > 0' % { sel: sel },
            'for': '%s' % c.alertFor,
            labels: { severity: 'critical' },
            annotations: {
              summary: 'SMART media errors (instance {{ $labels.instance }})',
              description: 'Disk controller detected media errors on {{ $labels.instance }} drive {{ $labels.device }}. VALUE = {{ $value }}. LABELS = {{ $labels }}',
            },
          },

          // ----------------------------------------------------------------
          // Wear-out
          // ----------------------------------------------------------------
          {
            alert: 'SmartWearoutIndicator',
            expr: 'smartctl_device_available_spare{%(sel)s} < smartctl_device_available_spare_threshold{%(sel)s}' % { sel: sel },
            'for': '%s' % c.alertFor,
            labels: { severity: 'critical' },
            annotations: {
              summary: 'SMART Wearout Indicator (instance {{ $labels.instance }})',
              description: 'Device is wearing out on {{ $labels.instance }} drive {{ $labels.device }}. VALUE = {{ $value }}. LABELS = {{ $labels }}',
            },
          },
          {
            alert: 'SmartDeviceSSDWearingOut',
            expr: 'smartctl_device_attribute{%(sel)s, attribute_name="Media_Wearout_Indicator", attribute_value_type="value"} < %(ssdWearoutThreshold)d' % { sel: sel, ssdWearoutThreshold: c.ssdWearoutThreshold },
            'for': '%s' % c.alertFor,
            labels: { severity: 'warning' },
            annotations: {
              summary: 'SSD is wearing out (instance {{ $labels.instance }})',
              description: 'Less than %(ssdWearoutThreshold)d%% lifetime remaining on SSD on {{ $labels.instance }} drive {{ $labels.device }}. VALUE = {{ $value }}. LABELS = {{ $labels }}' % c,
            },
          },

          // ----------------------------------------------------------------
          // Error logs
          // ----------------------------------------------------------------
          {
            alert: 'SmartErrorLogs',
            expr: 'smartctl_device_error_log_count{%(sel)s} > 0' % { sel: sel },
            'for': '%s' % c.alertFor,
            labels: { severity: 'critical' },
            annotations: {
              summary: 'SMART Warning Logs on {{ $labels.instance }}',
              description: 'Device logged warnings on {{ $labels.instance }} (drive {{ $labels.device }}). VALUE = {{ $value }}. LABELS = {{ $labels }}',
            },
          },
        ],
      },
    ],
  },
}
