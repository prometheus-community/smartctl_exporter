// Copyright 2022 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

// Device
type Device struct {
	Name  string
	Type  string
	Label string
}

func (d Device) String() string {
	return d.Name + ";" + d.Type + " (" + d.Label + ")"
}

// SMARTctlManagerCollector implements the Collector interface.
type SMARTctlManagerCollector struct {
	CollectPeriod         string
	CollectPeriodDuration time.Duration
	Devices               []Device

	logger *slog.Logger
	mutex  sync.Mutex
}

// Describe sends the super-set of all possible descriptors of metrics
func (i *SMARTctlManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricSmartctlVersion
	ch <- metricDeviceModel
	ch <- metricDeviceCount
	ch <- metricDeviceCapacityBlocks
	ch <- metricDeviceCapacityBytes
	ch <- metricDeviceTotalCapacityBytes
	ch <- metricDeviceBlockSize
	ch <- metricDeviceInterfaceSpeed
	ch <- metricDeviceAttribute
	ch <- metricDevicePowerOnSeconds
	ch <- metricDeviceRotationRate
	ch <- metricDeviceTemperature
	ch <- metricDevicePowerCycleCount
	ch <- metricDevicePercentageUsed
	ch <- metricDeviceAvailableSpare
	ch <- metricDeviceAvailableSpareThreshold
	ch <- metricDeviceCriticalWarning
	ch <- metricDeviceMediaErrors
	ch <- metricDeviceNumErrLogEntries
	ch <- metricDeviceBytesRead
	ch <- metricDeviceBytesWritten
	ch <- metricDeviceSmartStatus
	ch <- metricDeviceExitStatus
	ch <- metricDeviceState
	ch <- metricDeviceStatistics
	ch <- metricDeviceErrorLogCount
	ch <- metricDeviceSelfTestLogCount
	ch <- metricDeviceSelfTestLogErrorCount
	ch <- metricDeviceERCSeconds
	ch <- metricSCSIGrownDefectList
	ch <- metricReadErrorsCorrectedByRereadsRewrites
	ch <- metricReadErrorsCorrectedByEccFast
	ch <- metricReadErrorsCorrectedByEccDelayed
	ch <- metricReadTotalUncorrectedErrors
	ch <- metricWriteErrorsCorrectedByRereadsRewrites
	ch <- metricWriteErrorsCorrectedByEccFast
	ch <- metricWriteErrorsCorrectedByEccDelayed
	ch <- metricWriteTotalUncorrectedErrors
}

// Collect is called by the Prometheus registry when collecting metrics.
func (i *SMARTctlManagerCollector) Collect(ch chan<- prometheus.Metric) {
	info := NewSMARTctlInfo(ch)
	i.mutex.Lock()
	refreshAllDevices(i.logger, i.Devices)
	for _, device := range i.Devices {
		json := readData(i.logger, device)
		if json.Exists() {
			info.SetJSON(json)
			smart := NewSMARTctl(i.logger, json, ch, device)
			smart.Collect()
		}
	}
	ch <- prometheus.MustNewConstMetric(
		metricDeviceCount,
		prometheus.GaugeValue,
		float64(len(i.Devices)),
	)
	info.Collect()
	i.mutex.Unlock()
}

func (i *SMARTctlManagerCollector) RescanForDevices() {
	for {
		time.Sleep(*smartctlRescanInterval)
		i.logger.Info("Rescanning for devices")
		devices := scanDevices(i.logger)
		devices = append(devices, scanExtraDevices(i.logger, *smartctlExtraDevices)...)
		devices = buildDevicesFromFlag(devices)
		i.mutex.Lock()
		i.Devices = devices
		i.mutex.Unlock()
	}
}

var (
	smartctlPath = kingpin.Flag("smartctl.path",
		"The path to the smartctl binary",
	).Default("/usr/sbin/smartctl").String()
	smartctlInterval = kingpin.Flag("smartctl.interval",
		"The interval between smartctl polls",
	).Default("60s").Duration()
	smartctlRescanInterval = kingpin.Flag("smartctl.rescan",
		"The interval between rescanning for new/disappeared devices. If the interval is smaller than 1s no rescanning takes place. If any devices are configured with smartctl.device also no rescanning takes place.",
	).Default("10m").Duration()
	smartctlScan    = kingpin.Flag("smartctl.scan", "Enable scanning. This is a default if no devices are specified").Default("false").Bool()
	smartctlDevices = kingpin.Flag("smartctl.device",
		"The device to monitor. Device type can be specified after a semicolon, eg. '/dev/bus/0;megaraid,1' (repeatable)",
	).Strings()
	smartctlDeviceExclude = kingpin.Flag(
		"smartctl.device-exclude",
		"Regexp of devices to exclude from automatic scanning. (mutually exclusive to device-include)",
	).Default("").String()
	smartctlDeviceInclude = kingpin.Flag(
		"smartctl.device-include",
		"Regexp of devices to exclude from automatic scanning. (mutually exclusive to device-exclude)",
	).Default("").String()
	smartctlScanDeviceTypes = kingpin.Flag(
		"smartctl.scan-device-type",
		"Device type to use during automatic scan. Special by-id value forces predictable device names. (repeatable)",
	).Strings()
	smartctlExtraDevices = kingpin.Flag(
		"smartctl.extra-device",
		"Extra device scan specification in the format '<directory>;<name-regex>;<device-type>'. Scans the directory for entries matching the regex and monitors them with the given device type (repeatable). Example: '/dev/spdk;^nvme\\d+$;nvme'",
	).Strings()
	smartctlFakeData = kingpin.Flag("smartctl.fake-data",
		"The device to monitor (repeatable)",
	).Default("false").Hidden().Bool()
	smartctlPowerModeCheck = kingpin.Flag("smartctl.powermode-check",
		"Whether or not to check powermode before fetching data",
	).Default("standby").String()
)

// scanDevices uses smartctl to gather the list of available devices.
func scanDevices(logger *slog.Logger) []Device {
	filter := newDeviceFilter(*smartctlDeviceExclude, *smartctlDeviceInclude)

	json := readSMARTctlDevices(logger)
	scanDevices := json.Get("devices").Array()
	var scanDeviceResult []Device
	for _, d := range scanDevices {
		deviceName := d.Get("name").String()
		deviceType := d.Get("type").String()

		// SATA devices are reported as SCSI during scan - fallback to auto scraping
		if deviceType == "scsi" {
			deviceType = "auto"
		}

		deviceLabel := buildDeviceLabel(deviceName, deviceType)
		if filter.ignored(deviceLabel) {
			logger.Info("Ignoring device", "name", deviceLabel)
		} else {
			logger.Info("Found device", "name", deviceLabel)
			device := Device{
				Name:  deviceName,
				Type:  deviceType,
				Label: deviceLabel,
			}
			scanDeviceResult = append(scanDeviceResult, device)
		}
	}
	return scanDeviceResult
}

// extraDeviceSpec holds the parsed form of a --smartctl.extra-device value.
type extraDeviceSpec struct {
	directory  string
	nameRegexp *regexp.Regexp
	deviceType string
}

// parseExtraDeviceSpec parses "<directory>;<name-regex>;<device-type>".
func parseExtraDeviceSpec(s string) (extraDeviceSpec, error) {
	parts := strings.SplitN(s, ";", 3)
	if len(parts) != 3 {
		return extraDeviceSpec{}, fmt.Errorf("extra-device %q: expected format '<directory>;<name-regex>;<device-type>'", s)
	}
	re, err := regexp.Compile(parts[1])
	if err != nil {
		return extraDeviceSpec{}, fmt.Errorf("extra-device %q: invalid name regex: %w", s, err)
	}
	return extraDeviceSpec{directory: parts[0], nameRegexp: re, deviceType: parts[2]}, nil
}

// scanExtraDevices reads each extra-device spec, enumerates matching entries in the
// specified directory, applies the existing include/exclude filter, and returns Devices.
func scanExtraDevices(logger *slog.Logger, specs []string) []Device {
	filter := newDeviceFilter(*smartctlDeviceExclude, *smartctlDeviceInclude)
	var result []Device
	for _, raw := range specs {
		spec, err := parseExtraDeviceSpec(raw)
		if err != nil {
			logger.Error("Skipping invalid extra-device spec", "spec", raw, "err", err)
			continue
		}
		entries, err := os.ReadDir(spec.directory)
		if err != nil {
			logger.Error("Cannot read extra-device directory", "directory", spec.directory, "err", err)
			continue
		}
		for _, entry := range entries {
			if !spec.nameRegexp.MatchString(entry.Name()) {
				continue
			}
			deviceName := spec.directory + "/" + entry.Name()
			deviceLabel := buildDeviceLabel(deviceName, spec.deviceType)
			if filter.ignored(deviceLabel) {
				logger.Info("Ignoring extra device", "name", deviceLabel)
				continue
			}
			logger.Info("Found extra device", "name", deviceLabel)
			result = append(result, Device{
				Name:  deviceName,
				Type:  spec.deviceType,
				Label: deviceLabel,
			})
		}
	}
	return result
}

func buildDevicesFromFlag(devices []Device) []Device {
	// TODO: deduplication?
	for _, device := range *smartctlDevices {
		deviceName, deviceType, _ := strings.Cut(device, ";")
		if deviceType == "" {
			deviceType = "auto"
		}

		devices = append(devices, Device{
			Name:  deviceName,
			Type:  deviceType,
			Label: buildDeviceLabel(deviceName, deviceType),
		})
	}
	return devices
}

func validatePowerMode(mode string) error {
	switch strings.ToLower(mode) {
	case "never", "sleep", "standby", "idle":
		return nil
	default:
		return fmt.Errorf("invalid power mode: %s. Must be one of: never, sleep, standby, idle", mode)
	}
}

func main() {
	metricsPath := kingpin.Flag(
		"web.telemetry-path", "Path under which to expose metrics",
	).Default("/metrics").String()
	toolkitFlags := webflag.AddFlags(kingpin.CommandLine, ":9633")

	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("smartctl_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promslog.New(promslogConfig)

	if err := validatePowerMode(*smartctlPowerModeCheck); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Starting smartctl_exporter", "version", version.Info())
	logger.Info("Build context", "build_context", version.BuildContext())
	var devices []Device

	if len(*smartctlDevices) == 0 {
		*smartctlScan = true
	}

	if *smartctlScan {
		devices = scanDevices(logger)
		logger.Info("Number of devices found", "count", len(devices))
	}

	if len(*smartctlExtraDevices) > 0 {
		logger.Info("Extra device specs specified", "specs", strings.Join(*smartctlExtraDevices, ", "))
		extraDevices := scanExtraDevices(logger, *smartctlExtraDevices)
		devices = append(devices, extraDevices...)
		logger.Info("Extra devices found", "count", len(extraDevices))
	}

	if len(*smartctlDevices) > 0 {
		logger.Info("Devices specified", "devices", strings.Join(*smartctlDevices, ", "))
		devices = buildDevicesFromFlag(devices)
		logger.Info("Devices filtered", "count", len(devices))
	}

	collector := SMARTctlManagerCollector{
		Devices: devices,
		logger:  logger,
	}

	if (*smartctlScan || len(*smartctlExtraDevices) > 0) && *smartctlRescanInterval >= 1*time.Second {
		logger.Info("Start background scan process")
		logger.Info("Rescanning for devices every", "rescanInterval", *smartctlRescanInterval)
		go collector.RescanForDevices()
	}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(versioncollector.NewCollector("smartctl_exporter"))
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	prometheus.WrapRegistererWithPrefix("", reg).MustRegister(&collector)

	http.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	if *metricsPath != "/" && *metricsPath != "" {
		landingConfig := web.LandingConfig{
			Name:        "smartctl_exporter",
			Description: "Prometheus Exporter for S.M.A.R.T. devices",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error("error creating landing page", "err", err)
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
		logger.Error("error running HTTP server", "err", err)
		os.Exit(1)
	}
}
