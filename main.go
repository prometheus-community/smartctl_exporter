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
	"runtime"
	"strconv"
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
	prometheus.DescribeByCollect(i, ch)
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
			smart := NewSMARTctl(i.logger, json, ch)
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
		devices = buildDevicesFromFlag(devices)
		i.mutex.Lock()
		i.Devices = devices
		i.mutex.Unlock()
	}
}

var (
	configFile *string

	smartctlPath            *string
	smartctlInterval        *time.Duration
	smartctlRescanInterval  *time.Duration
	smartctlScan            *bool
	smartctlDevices         *[]string
	smartctlDeviceExclude   *string
	smartctlDeviceInclude   *string
	smartctlScanDeviceTypes *[]string
	smartctlFakeData        *bool
	smartctlPowerModeCheck  *string
)

func defaultSmartctlPath() string {
	if runtime.GOOS == "windows" {
		return "smartctl.exe"
	}
	return "/usr/sbin/smartctl"
}

func initFlags(cfg Config, configPath string) (*string, *web.FlagConfig) {
	configFile = kingpin.Flag("config.file", "Path to configuration file.").Default(configPath).String()

	smartctlPath = kingpin.Flag("smartctl.path",
		"The path to the smartctl binary",
	).Default(stringValueOrDefault(cfg.SmartctlPath, defaultSmartctlPath())).String()
	smartctlInterval = kingpin.Flag("smartctl.interval",
		"The interval between smartctl polls",
	).Default(stringValueOrDefault(cfg.SmartctlInterval, "60s")).Duration()
	smartctlRescanInterval = kingpin.Flag("smartctl.rescan",
		"The interval between rescanning for new/disappeared devices. If the interval is smaller than 1s no rescanning takes place. If any devices are configured with smartctl.device also no rescanning takes place.",
	).Default(stringValueOrDefault(cfg.SmartctlRescan, "10m")).Duration()
	smartctlScan = kingpin.Flag("smartctl.scan", "Enable scanning. This is a default if no devices are specified").Default(strconv.FormatBool(cfg.SmartctlScan != nil && *cfg.SmartctlScan)).Bool()
	devicesFlag := kingpin.Flag("smartctl.device",
		"The device to monitor. Device type can be specified after a semicolon, eg. '/dev/bus/0;megaraid,1' (repeatable)",
	)
	if len(cfg.SmartctlDevices) > 0 {
		devicesFlag.Default(cfg.SmartctlDevices...)
	}
	smartctlDevices = devicesFlag.Strings()
	smartctlDeviceExclude = kingpin.Flag(
		"smartctl.device-exclude",
		"Regexp of devices to exclude from automatic scanning. (mutually exclusive to device-include)",
	).Default(stringValueOrDefault(cfg.SmartctlDeviceExclude, "")).String()
	smartctlDeviceInclude = kingpin.Flag(
		"smartctl.device-include",
		"Regexp of devices to exclude from automatic scanning. (mutually exclusive to device-exclude)",
	).Default(stringValueOrDefault(cfg.SmartctlDeviceInclude, "")).String()
	scanTypesFlag := kingpin.Flag(
		"smartctl.scan-device-type",
		"Device type to use during automatic scan. Special by-id value forces predictable device names. (repeatable)",
	)
	if len(cfg.SmartctlScanDeviceTypes) > 0 {
		scanTypesFlag.Default(cfg.SmartctlScanDeviceTypes...)
	}
	smartctlScanDeviceTypes = scanTypesFlag.Strings()
	smartctlFakeData = kingpin.Flag("smartctl.fake-data",
		"The device to monitor (repeatable)",
	).Default("false").Hidden().Bool()
	smartctlPowerModeCheck = kingpin.Flag("smartctl.powermode-check",
		"Whether or not to check powermode before fetching data",
	).Default(stringValueOrDefault(cfg.SmartctlPowerModeCheck, "standby")).String()

	metricsPath := kingpin.Flag(
		"web.telemetry-path", "Path under which to expose metrics",
	).Default(stringValueOrDefault(cfg.WebTelemetryPath, "/metrics")).String()
	toolkitFlags := webflag.AddFlags(kingpin.CommandLine, stringValueOrDefault(cfg.WebListenAddress, ":9633"))
	return metricsPath, toolkitFlags
}

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
	configPath := findConfigFile(os.Args[1:])
	cfg, err := loadConfigFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	metricsPath, toolkitFlags := initFlags(cfg, configPath)

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

	if len(*smartctlDevices) > 0 {
		logger.Info("Devices specified", "devices", strings.Join(*smartctlDevices, ", "))
		devices = buildDevicesFromFlag(devices)
		logger.Info("Devices filtered", "count", len(devices))
	}

	collector := SMARTctlManagerCollector{
		Devices: devices,
		logger:  logger,
	}

	if *smartctlScan && *smartctlRescanInterval >= 1*time.Second {
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
