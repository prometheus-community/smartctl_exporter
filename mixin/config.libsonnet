{
  _config+:: {
    // Select the metrics coming from the smartctl_exporter. Adjust to match
    // your scrape job, e.g. 'job="smartctl-exporter"'.
    smartctlExporterSelector: 'job="smartctl"',

    // Exclude virtual disks (e.g. CSI / iSCSI) from dashboards and rules. By default, filter out virtual disks created by LongHorn in K8s.
    // Empty string means no filtering.
    deviceModelExclusionSelector: 'model_name!="IET VIRTUAL-DISK"',

    // Absolute temperature thresholds (°C) used by the temperature alerts
    // that do not rely on the drive trip value.
    temperatureWarningThreshold: 60,
    temperatureCriticalThreshold: 70,

    // Fraction of the drive trip value at which the "nearing trip"
    // warning alert fires. 0.80 means 80% of the trip value is "nearing"
    temperatureNearingTripRatio: 0.80,

    // Remaining SSD lifetime (%) below which SmartDeviceSSDWearingOut alert fires.
    ssdWearoutThreshold: 10,

    // Look-back windows for temperature smoothing before evaluating the
    // absolute thresholds.
    temperatureAvgWindow: '5m',
    temperatureMaxWindow: '10m',

    // How long conditions must hold before alerting.
    alertFor: '10m',

    // Prometheus rule group evaluation interval.
    ruleEvaluationInterval: '10m',

    dashboardNamePrefix: 'smartctl / ',
    dashboardTags: ['smartctl-mixin'],
  },
}
