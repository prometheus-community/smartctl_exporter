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
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// SMARTctlManagerCollector implements the Collector interface.
type SMARTctlManagerCollector struct {
	CollectPeriod         string
	CollectPeriodDuration time.Duration
	Devices               []string

	logger log.Logger
}

// Describe sends the super-set of all possible descriptors of metrics
func (i SMARTctlManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(i, ch)
}

// Collect is called by the Prometheus registry when collecting metrics.
func (i SMARTctlManagerCollector) Collect(ch chan<- prometheus.Metric) {
	info := NewSMARTctlInfo(ch)
	for _, device := range i.Devices {
		json := readData(i.logger, device)
		if json.Exists() {
			info.SetJSON(json)
			smart := NewSMARTctl(i.logger, json, ch)
			smart.Collect()
		}
	}
	info.Collect()
}

var (
	smartctlPath = kingpin.Flag("smartctl.path",
		"The path to the smartctl binary",
	).Default("/usr/sbin/smartctl").String()
	smartctlInterval = kingpin.Flag("smartctl.interval",
		"The interval between smarctl polls",
	).Default("60s").Duration()
	smartctlDevices = kingpin.Flag("smartctl.device",
		"The device to monitor (repeatable)",
	).Strings()
	smartctlDeviceExclude = kingpin.Flag(
		"smartctl.device-exclude",
		"Regexp of devices to exclude.",
	).Default("").String()
	smartctlFakeData = kingpin.Flag("smartctl.fake-data",
		"The device to monitor (repeatable)",
	).Default("false").Hidden().Bool()
)

type deviceFilter struct {
	ignorePattern *regexp.Regexp
}

func newDeviceFilter(ignoredPattern string) (f deviceFilter) {
	if ignoredPattern != "" {
		f.ignorePattern = regexp.MustCompile(ignoredPattern)
	}
	return
}

func (f *deviceFilter) ignored(name string) bool {
	return (f.ignorePattern != nil && f.ignorePattern.MatchString(name))
}

func main() {
	metricsPath := kingpin.Flag(
		"web.telemetry-path", "Path under which to expose metrics",
	).Default("/metrics").String()
	toolkitFlags := webflag.AddFlags(kingpin.CommandLine, ":9633")

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("smartctl_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting smartctl_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	// Scan the host devices
	json := readSMARTctlDevices(logger)
	scanDevices := json.Get("devices").Array()
	scanDevicesSet := make(map[string]bool)
	var scanDeviceNames []string
	for _, d := range scanDevices {
		deviceName := d.Get("name").String()
		level.Debug(logger).Log("msg", "Found device", "name", deviceName)
		scanDevicesSet[deviceName] = true
		scanDeviceNames = append(scanDeviceNames, deviceName)
	}

	// Read the configuration and verify that it is available
	devices := *smartctlDevices
	var readDeviceNames []string
	for _, device := range devices {
		if _, ok := scanDevicesSet[device]; ok {
			readDeviceNames = append(readDeviceNames, device)
		} else {
			level.Warn(logger).Log("msg", "Device unavailable", "name", device)
		}
	}

	if len(readDeviceNames) > 0 {
		devices = readDeviceNames
	} else {
		level.Info(logger).Log("msg", "No devices specified, trying to load them automatically")
		devices = scanDeviceNames
	}
	filter := newDeviceFilter(*smartctlDeviceExclude)
	for i, device := range devices {
		if filter.ignored(device) {
			devices = append(devices[:i], devices[i+1:]...)
		}
	}

	if len(devices) == 0 {
		level.Error(logger).Log("msg", "No devices found")
		os.Exit(1)
	}

	collector := SMARTctlManagerCollector{
		Devices: devices,
		logger:  logger,
	}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	prometheus.WrapRegistererWithPrefix("", reg).MustRegister(collector)

	http.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
      <head><title>Smartctl Exporter</title></head>
      <body>
      <h1>Smartctl Exporter</h1>
      <p><a href="` + *metricsPath + `">Metrics</a></p>
      </body>
      </html>`))
		if err != nil {
			level.Error(logger).Log("msg", "Couldn't write response", "err", err)
		}
	})

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
