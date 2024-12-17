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
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

// Device
type Device struct {
	Name      string `json:"name"`
	Info_Name string `json:"info_name"`
	Type      string `json:"type"`
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
	smartctlDevices = kingpin.Flag("smartctl.device",
		"The device to monitor (repeatable)",
	).Strings()
	smartctlDeviceExclude = kingpin.Flag(
		"smartctl.device-exclude",
		"Regexp of devices to exclude from automatic scanning. (mutually exclusive to device-include)",
	).Default("").String()
	smartctlDeviceInclude = kingpin.Flag(
		"smartctl.device-include",
		"Regexp of devices to exclude from automatic scanning. (mutually exclusive to device-exclude)",
	).Default("").String()
	smartctlFakeData = kingpin.Flag("smartctl.fake-data",
		"The device to monitor (repeatable)",
	).Default("false").Hidden().Bool()
)

// scanDevices uses smartctl to gather the list of available devices.
func scanDevices(logger *slog.Logger) []Device {
	filter := newDeviceFilter(*smartctlDeviceExclude, *smartctlDeviceInclude)

	json := readSMARTctlDevices(logger)
	scanDevices := json.Get("devices").Array()
	var scanDeviceResult []Device
	for _, d := range scanDevices {
		deviceName := extractDiskName(strings.TrimSpace(d.Get("info_name").String()))
		if filter.ignored(deviceName) {
			logger.Info("Ignoring device", "name", deviceName)
		} else if !d.Get("smart_support.available").Bool() {
			logger.Info("Ignoring device", "name", deviceName, "reason:", "SMART not available")
		} else {
			logger.Info("Found device", "name", deviceName)
			device := Device{
				Name:      d.Get("name").String(),
				Info_Name: deviceName,
				Type:      d.Get("type").String(),
			}
			scanDeviceResult = append(scanDeviceResult, device)
		}
	}
	return scanDeviceResult
}

func filterDevices(logger *slog.Logger, devices []Device, filters []string) []Device {
	var filtered []Device
	for _, d := range devices {
		for _, filter := range filters {
			logger.Debug("filterDevices", "device", d.Info_Name, "filter", filter)
			if strings.Contains(d.Info_Name, filter) {
				filtered = append(filtered, d)
				break
			}
		}
	}
	return filtered
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

	logger.Info("Starting smartctl_exporter", "version", version.Info())
	logger.Info("Build context", "build_context", version.BuildContext())

	var devices []Device
	devices = scanDevices(logger)
	logger.Info("Number of devices found", "count", len(devices))
	if len(*smartctlDevices) > 0 {
		logger.Info("Devices specified", "devices", strings.Join(*smartctlDevices, ", "))
		devices = filterDevices(logger, devices, *smartctlDevices)
		logger.Info("Devices filtered", "count", len(devices))
	}

	collector := SMARTctlManagerCollector{
		Devices: devices,
		logger:  logger,
	}

	if *smartctlRescanInterval >= 1*time.Second {
		logger.Info("Start background scan process")
		logger.Info("Rescanning for devices every", "rescanInterval", *smartctlRescanInterval)
		go collector.RescanForDevices()
	}

	reg := prometheus.NewPedanticRegistry()
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
