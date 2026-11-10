{
  local c = $._config,
  local sel = c.smartctlExporterSelector,
  local excl = c.deviceModelExclusionSelector,
  local ds = { type: 'prometheus', uid: '${DS_PROMETHEUS}' },

  local prom(expr, legend, instant=false, ref='A') = {
    datasource: ds,
    editorMode: 'code',
    exemplar: false,
    expr: expr,
    instant: instant,
    legendFormat: legend,
    range: !instant,
    refId: ref,
  },

  local step(color, value=null) =
    if value == null then { color: color } else { color: color, value: value },

  local tsDefaults(unit, thresholds, max=null) = {
    unit: unit,
    color: { mode: 'palette-classic' },
    max: if max == null then null else max,
    thresholds: { mode: 'absolute', steps: thresholds },
    custom: {
      axisPlacement: 'auto',
      drawStyle: 'line',
      fillOpacity: 10,
      lineInterpolation: 'smooth',
      lineWidth: 1,
      pointSize: 5,
      showPoints: 'auto',
      spanNulls: false,
      stacking: { group: 'A', mode: 'none' },
      thresholdsStyle: { mode: 'off' },
    },
  },

  local tsOptions = {
    legend: { calcs: ['lastNotNull', 'max'], displayMode: 'table', placement: 'bottom', showLegend: true },
    tooltip: { mode: 'multi', sort: 'desc' },
  },

  local statDefaults(thresholds) = {
    unit: 'short',
    thresholds: { mode: 'absolute', steps: thresholds },
  },

  local tempQuery(instant=false, legend='{{instance}}/{{device}}', ref='Temperature') =
    if excl == '' then
      prom('sum by (instance, device)(smartctl_device_temperature{%(sel)s, instance=~"$instance", temperature_type="current"})' % { sel: sel }, legend, instant, ref)
    else
      prom('sum by (instance, device)(smartctl_device_temperature{%(sel)s, instance=~"$instance", temperature_type="current"}) and on(instance, device) smartctl_device{%(excl)s}' % { sel: sel, excl: excl }, legend, instant, ref),

  local tempMedianQuery =
    if excl == '' then
      prom('quantile_over_time(0.5, sum by (instance, device) (smartctl_device_temperature{%(sel)s, instance=~"$instance", temperature_type="current"})[4w:1h])' % { sel: sel }, '{{instance}}/{{device}} median', true, 'Median temperature')
    else
      prom('quantile_over_time(0.5, sum by (instance, device) (smartctl_device_temperature{%(sel)s, instance=~"$instance", temperature_type="current"} and on(instance, device) smartctl_device{%(excl)s})[4w:1h])' % { sel: sel, excl: excl }, '{{instance}}/{{device}} median', true, 'Median temperature'),

  local devicesQuery =
    if excl == '' then
      prom('smartctl_device{%(sel)s, instance=~"$instance"}' % { sel: sel }, '__auto', true, 'Devices') + { format: 'table' }
    else
      prom('smartctl_device{%(sel)s, instance=~"$instance", %(excl)s}' % { sel: sel, excl: excl }, '__auto', true, 'Devices') + { format: 'table' },

  grafanaDashboards+:: {
    'smartctl-overview.json': {
      uid: 'smartctl-overview',
      title: 'SMART Disk Health Monitoring',
      tags: c.dashboardTags,
      timezone: 'browser',
      schemaVersion: 39,
      version: 1,
      refresh: '',
      time: { from: 'now-12h', to: 'now' },
      timepicker: {
        refresh_intervals: ['5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h', '2h', '1d'],
      },
      templating: {
        list: [
          {
            name: 'DS_PROMETHEUS',
            type: 'datasource',
            query: 'prometheus',
            label: 'Prometheus datasource',
            current: { text: 'Prometheus', value: 'Prometheus' },
            hide: 0,
          },
          {
            name: 'instance',
            type: 'query',
            datasource: ds,
            label: 'Instance',
            query: 'label_values(smartctl_device_smart_status{%(sel)s}, instance)' % { sel: sel },
            refresh: 1,
            sort: 1,
            multi: true,
            includeAll: true,
            current: { selected: false, text: 'All', value: '$__all' },
            options: [],
          },
        ],
      },
      panels: [
        // ----------------------------------------------------------------
        // Header
        // ----------------------------------------------------------------
        {
          id: 31,
          type: 'text',
          gridPos: { h: 4, w: 24, x: 0, y: 0 },
          transparent: true,
          options: {
            mode: 'html',
            content: '<div style="background: linear-gradient(135deg, #0d1117 0%, #1a1f2e 60%, #0f172a 100%); padding: 18px 28px; display: flex; align-items: center; gap: 22px; border-radius: 6px; border-left: 4px solid #3b82f6;">\n  <div>\n    <div style="color: #e2e8f0; margin: 0; font-size: 26px; font-weight: 700;">S.M.A.R.T. Disk Health Monitoring</div>\n    <div style="color: #64748b; margin: 5px 0 0; font-size: 13px;">Self-Monitoring, Analysis and Reporting Technology &mdash; real-time disk diagnostics &amp; predictive failure analysis</div>\n  </div>\n</div>',
          },
        },

        // ----------------------------------------------------------------
        // Overview row
        // ----------------------------------------------------------------
        {
          id: 100,
          type: 'row',
          title: 'Overview',
          gridPos: { h: 1, w: 24, x: 0, y: 4 },
          collapsed: false,
          panels: [],
        },

        // SMART probe
        {
          id: 19,
          type: 'stat',
          title: 'SMART probe',
          description: 'Exit status from smartctl command for all devices selected. Indicates if metrics are collected properly.',
          datasource: ds,
          gridPos: { h: 3, w: 4, x: 0, y: 5 },
          targets: [
            prom('max(smartctl_device_smartctl_exit_status{%(sel)s, instance=~"$instance"} > 0 or vector(0))' % { sel: sel }, '__auto', false, 'Probe status'),
          ],
          fieldConfig: {
            defaults: {
              unit: 'short',
              mappings: [
                { type: 'value', options: { '0': { text: 'Running', index: 0 } } },
                { type: 'range', options: { from: 1, to: null, result: { text: 'smartctl errors', index: 1 } } },
              ],
              thresholds: { mode: 'absolute', steps: [step('blue'), step('red', 1)] },
            },
            overrides: [],
          },
          options: {
            colorMode: 'background',
            graphMode: 'none',
            reduceOptions: { calcs: ['max'], fields: '', values: false },
            textMode: 'value',
            wideLayout: true,
          },
        },

        // Disk power-on time
        {
          id: 7,
          type: 'stat',
          title: 'Disk power-on time',
          description: 'Number of seconds each disk has been powered on',
          datasource: ds,
          gridPos: { h: 7, w: 8, x: 4, y: 5 },
          targets: [
            prom('max by (instance, device)(smartctl_device_power_on_seconds{%(sel)s, instance=~"$instance"})' % { sel: sel }, '{{instance}} - {{device}}', true, 'On-time'),
          ],
          fieldConfig: {
            defaults: {
              unit: 's',
              decimals: 1,
              color: { fixedColor: 'text', mode: 'fixed' },
              thresholds: { mode: 'absolute', steps: [step('text')] },
            },
            overrides: [],
          },
          options: {
            colorMode: 'value',
            graphMode: 'none',
            reduceOptions: { calcs: ['lastNotNull'], fields: '', values: false },
            textMode: 'auto',
            wideLayout: true,
          },
        },

        // Disk temperature
        {
          id: 6,
          type: 'timeseries',
          title: 'Disk temperature',
          description: 'Disk temperatures over time for each device',
          datasource: ds,
          gridPos: { h: 9, w: 12, x: 12, y: 5 },
          targets: [
            tempQuery(false, '{{instance}}/{{device}}', 'Temperature'),
          ],
          fieldConfig: {
            defaults: tsDefaults('celsius', [step('blue'), step('yellow', 50), step('red', 60)], 70) + {
              custom+: { thresholdsStyle+: { mode: 'dashed' } },
            },
            overrides: [],
          },
          options: {
            legend: {
              calcs: ['lastNotNull', 'median'],
              displayMode: 'list',
              placement: 'bottom',
              showLegend: true,
            },
            tooltip: { mode: 'multi', sort: 'desc' },
          },
        },

        // Failed disks count
        {
          id: 20,
          type: 'stat',
          title: 'Failed disks count',
          description: 'Number of disks with failing SMART status. Indicates disk death.',
          datasource: ds,
          gridPos: { h: 4, w: 4, x: 0, y: 8 },
          targets: [
            prom('count by (instance)(smartctl_device{%(sel)s, instance=~"$instance"}) - count by (instance)(smartctl_device_smart_status{%(sel)s, instance=~"$instance"} == 1)' % { sel: sel }, '__auto', true, 'Failing disks'),
          ],
          fieldConfig: {
            defaults: statDefaults([step('transparent'), step('red', 1)]),
            overrides: [],
          },
          options: {
            colorMode: 'background',
            graphMode: 'area',
            reduceOptions: { calcs: ['lastNotNull'], fields: '', values: false },
            textMode: 'value_and_name',
            wideLayout: true,
          },
        },

        // Devices in scope
        {
          id: 30,
          type: 'table',
          title: 'Devices in scope',
          datasource: ds,
          gridPos: { h: 9, w: 12, x: 0, y: 12 },
          targets: [devicesQuery],
          transformations: [
            {
              id: 'organize',
              options: {
                excludeByName: {
                  Time: true,
                  Value: true,
                  __name__: true,
                  ata_additional_product_id: true,
                  ata_version: true,
                  container: true,
                  firmware_version: true,
                  interface: true,
                  job: true,
                  k8s_cluster_name: true,
                  model_family: true,
                  namespace: true,
                  pod: true,
                  protocol: true,
                  sata_version: true,
                },
                indexByName: {
                  Time: 0,
                  Value: 15,
                  __name__: 1,
                  ata_additional_product_id: 2,
                  ata_version: 3,
                  device: 5,
                  firmware_version: 6,
                  form_factor: 7,
                  instance: 4,
                  interface: 8,
                  job: 9,
                  model_family: 10,
                  model_name: 11,
                  protocol: 12,
                  sata_version: 13,
                  serial_number: 14,
                },
              },
            },
          ],
          fieldConfig: {
            defaults: {
              color: { mode: 'thresholds' },
              thresholds: { mode: 'absolute', steps: [step('green'), step('red', 80)] },
              custom: { align: 'auto', filterable: true, cellOptions: { type: 'auto' } },
            },
            overrides: [],
          },
          options: { showHeader: true, cellHeight: 'sm', enablePagination: true },
        },

        // Prefailure checks below thresholds
        {
          id: 25,
          type: 'table',
          title: 'Prefailure checks below thresholds',
          description: 'Detailed prefailured checks below thresholds. Indicates the disks are wearing out and need to be replaced with no immediate failure.',
          datasource: ds,
          gridPos: { h: 7, w: 12, x: 12, y: 14 },
          targets: [
            prom('max by(instance, device, attribute_name)(smartctl_device_attribute{%(sel)s, instance=~"$instance", attribute_flags_short=~"P.....", attribute_value_type="value"}) <= max by(instance, device, attribute_name)(smartctl_device_attribute{%(sel)s, instance=~"$instance", attribute_flags_short=~"P.....", attribute_value_type="thresh"})' % { sel: sel }, '__auto', true, 'Prefailure checks') + { format: 'table' },
          ],
          transformations: [
            {
              id: 'organize',
              options: {
                excludeByName: { Time: true, Value: true },
                indexByName: { Time: 0, Value: 4, attribute_name: 3, device: 2, instance: 1 },
                renameByName: { attribute_name: '' },
              },
            },
          ],
          fieldConfig: {
            defaults: {
              color: { fixedColor: 'yellow', mode: 'fixed' },
              thresholds: { mode: 'absolute', steps: [step('transparent')] },
              custom: { align: 'auto', filterable: true, cellOptions: { type: 'auto' } },
            },
            overrides: [
              {
                matcher: { id: 'byName', options: 'attribute_name' },
                properties: [{ id: 'custom.cellOptions', value: { type: 'pill' } }],
              },
            ],
          },
          options: { showHeader: true, cellHeight: 'sm', enablePagination: true, frameIndex: 1, sortBy: [{ desc: false, displayName: 'attribute_id' }] },
        },

        // ----------------------------------------------------------------
        // HDD / SSD row
        // ----------------------------------------------------------------
        {
          id: 101,
          type: 'row',
          title: 'HDD / SSD',
          gridPos: { h: 1, w: 24, x: 0, y: 21 },
          collapsed: false,
          panels: [],
        },

        // Failing HDDs
        {
          id: 26,
          type: 'table',
          title: 'Failing HDDs',
          description: 'Disks with any error logs, unrecoverable/reallocated/pending sectors on HDD. Could indicate imminent failure or data corruption.',
          datasource: ds,
          gridPos: { h: 8, w: 6, x: 0, y: 22 },
          targets: [
            prom('sum by (instance, device)(max_over_time(smartctl_device_attribute{%(sel)s, instance=~"$instance", attribute_name=~"Offline_Uncorrectable|Current_Pending_Sector|Reallocated_Sector_Ct", attribute_value_type="raw"}[$__range])) > 0' % { sel: sel }, '{{instance}} - {{device}}', true, 'Failing HDD') + { format: 'table' },
          ],
          transformations: [
            {
              id: 'organize',
              options: {
                excludeByName: { Time: true, Value: true },
                renameByName: { instance: 'Instance', device: 'Device' },
              },
            },
          ],
          fieldConfig: {
            defaults: {
              color: { mode: 'thresholds' },
              thresholds: { mode: 'absolute', steps: [step('transparent'), step('red', 1)] },
              custom: { align: 'auto', filterable: true, cellOptions: { type: 'auto' } },
            },
            overrides: [],
          },
          options: { showHeader: true, cellHeight: 'sm', enablePagination: true },
        },

        // SSD wearout indicator
        {
          id: 29,
          type: 'gauge',
          title: 'SSD wearout indicator',
          description: 'Remaining life measured for SSDs',
          datasource: ds,
          gridPos: { h: 8, w: 6, x: 6, y: 22 },
          targets: [
            prom('smartctl_device_attribute{%(sel)s, instance=~"$instance", attribute_name="Media_Wearout_Indicator", attribute_value_type="value"}' % { sel: sel }, '{{instance}} - {{device}}', true, 'SSD wearout'),
          ],
          fieldConfig: {
            defaults: {
              unit: 'percent',
              min: 0,
              max: 100,
              color: { mode: 'thresholds' },
              thresholds: { mode: 'absolute', steps: [step('red'), step('yellow', 10), step('blue', 30)] },
            },
            overrides: [],
          },
          options: {
            orientation: 'auto',
            reduceOptions: { calcs: ['lastNotNull'], fields: '', values: false },
            showThresholdLabels: false,
            showThresholdMarkers: true,
            textMode: 'auto',
          },
        },

        // Offline uncorrectable sectors
        {
          id: 12,
          type: 'stat',
          title: 'Offline uncorrectable sectors',
          description: 'Uncorrectable sectors. Any value > 0 indicates data loss or corruption risk',
          datasource: ds,
          gridPos: { h: 8, w: 6, x: 12, y: 22 },
          targets: [
            prom('smartctl_device_attribute{%(sel)s, instance=~"$instance", attribute_name="Offline_Uncorrectable", attribute_value_type="raw"}' % { sel: sel }, '{{instance}} - {{device}}', false, 'Offline sectors'),
          ],
          fieldConfig: {
            defaults: statDefaults([step('text'), step('red', 1)]),
            overrides: [],
          },
          options: {
            colorMode: 'background',
            graphMode: 'none',
            reduceOptions: { calcs: ['lastNotNull'], fields: '', values: false },
            textMode: 'auto',
            wideLayout: true,
          },
        },

        // Current pending sectors
        {
          id: 11,
          type: 'stat',
          title: 'Current pending sectors',
          description: 'Sectors waiting to be reallocated. Any value > 0 is concerning and may indicate imminent disk failure',
          datasource: ds,
          gridPos: { h: 8, w: 6, x: 18, y: 22 },
          targets: [
            prom('smartctl_device_attribute{%(sel)s, instance=~"$instance", attribute_name="Current_Pending_Sector", attribute_value_type="raw"}' % { sel: sel }, '{{instance}} - {{device}}', false, 'Pending sectors'),
          ],
          fieldConfig: {
            defaults: statDefaults([step('text'), step('red', 1)]),
            overrides: [],
          },
          options: {
            colorMode: 'background',
            graphMode: 'none',
            reduceOptions: { calcs: ['lastNotNull'], fields: '', values: false },
            textMode: 'auto',
            wideLayout: true,
          },
        },

        // Reallocated sectors
        {
          id: 10,
          type: 'timeseries',
          title: 'Reallocated sectors',
          description: 'Count of reallocated sectors. Any value > 0 indicates disk is wearing out. Steep increase indicates potential imminent failure.',
          datasource: ds,
          gridPos: { h: 8, w: 12, x: 0, y: 30 },
          targets: [
            prom('sum by (instance, device)(smartctl_device_attribute{%(sel)s, instance=~"$instance", attribute_name="Reallocated_Sector_Ct", attribute_value_type="raw"})' % { sel: sel }, '{{instance}} - {{device}}', false, 'Reallocated sectors'),
          ],
          fieldConfig: {
            defaults: tsDefaults('short', [step('green'), step('yellow', 5), step('red', 50)]),
            overrides: [],
          },
          options: tsOptions,
        },

        // Error log count
        {
          id: 13,
          type: 'timeseries',
          title: 'Error log count',
          description: 'Count of errors logged by the disk controller',
          datasource: ds,
          gridPos: { h: 8, w: 12, x: 12, y: 30 },
          targets: [
            prom('sum by (instance, device)(smartctl_device_error_log_count{%(sel)s, instance=~"$instance"})' % { sel: sel }, '{{instance}} - {{device}}', false, 'Errors'),
          ],
          fieldConfig: {
            defaults: tsDefaults('short', [step('green'), step('yellow', 1), step('red', 10)]),
            overrides: [],
          },
          options: tsOptions,
        },

        // ----------------------------------------------------------------
        // NVMe row
        // ----------------------------------------------------------------
        {
          id: 102,
          type: 'row',
          title: 'NVMe',
          gridPos: { h: 1, w: 24, x: 0, y: 38 },
          collapsed: false,
          panels: [],
        },

        // NVMe percentage used
        {
          id: 15,
          type: 'gauge',
          title: 'NVMe percentage used',
          description: 'Percentage of NVMe device life used. Values approaching 100% indicate end of life',
          datasource: ds,
          gridPos: { h: 8, w: 6, x: 0, y: 39 },
          targets: [
            prom('smartctl_device_percentage_used{%(sel)s, instance=~"$instance"}' % { sel: sel }, '{{instance}} - {{device}}', true, 'NVMe life'),
          ],
          fieldConfig: {
            defaults: {
              unit: 'percent',
              min: 0,
              max: 100,
              color: { mode: 'thresholds' },
              thresholds: { mode: 'absolute', steps: [step('blue'), step('yellow', 80), step('red', 95)] },
            },
            overrides: [],
          },
          options: {
            orientation: 'auto',
            reduceOptions: { calcs: ['lastNotNull'], fields: '', values: false },
            showThresholdLabels: false,
            showThresholdMarkers: true,
            textMode: 'auto',
          },
        },

        // NVMe media errors
        {
          id: 14,
          type: 'timeseries',
          title: 'NVMe media errors',
          description: 'Media errors for NVMe devices. Any value > 0 indicates potential issues',
          datasource: ds,
          gridPos: { h: 8, w: 12, x: 6, y: 39 },
          targets: [
            prom('sum by (instance, device)(smartctl_device_media_errors{%(sel)s, instance=~"$instance"})' % { sel: sel }, '{{instance}} - {{device}}', false, 'NVme errors'),
          ],
          fieldConfig: {
            defaults: tsDefaults('short', [step('blue'), step('red', 1)]),
            overrides: [],
          },
          options: tsOptions,
        },
      ],
    },
  },
}
