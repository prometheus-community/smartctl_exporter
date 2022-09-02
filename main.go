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
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	options Options
	logger  Logger
)

// SMARTctlManagerCollector implements the Collector interface.
type SMARTctlManagerCollector struct {
}

// Describe sends the super-set of all possible descriptors of metrics
func (i SMARTctlManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(i, ch)
}

// Collect is called by the Prometheus registry when collecting metrics.
func (i SMARTctlManagerCollector) Collect(ch chan<- prometheus.Metric) {
	info := NewSMARTctlInfo(ch)
	for _, device := range options.SMARTctl.Devices {
		if json, err := readData(device); err == nil {
			info.SetJSON(json)
			smart := NewSMARTctl(json, ch)
			smart.Collect()
		} else {
			logger.Error(err.Error())
		}
	}
	info.Collect()
}

func init() {
	options = loadOptions()

	json := readSMARTctlDevices()
	devices := json.Get("devices").Array()
	deviceSet := make(map[string]bool)
	for _, d := range devices {
		device := d.Get("name").String()
		logger.Debug("Found device: %s", device)
		deviceSet[device] = true
	}

	var deviceList []string
	for _, device := range options.SMARTctl.Devices {
		if _, ok := deviceSet[device]; ok {
			deviceList = append(deviceList, device)
		} else {
			logger.Debug("Device %s unavialable", device)
		}
	}

	if len(options.SMARTctl.Devices) == 0 {
		logger.Debug("No devices specified, trying to load them automatically")
		options.SMARTctl.Devices = deviceList
	}
}

func main() {
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	prometheus.WrapRegistererWithPrefix("", reg).MustRegister(SMARTctlManagerCollector{})

	logger.Info("Starting on %s%s", options.SMARTctl.BindTo, options.SMARTctl.URLPath)
	http.Handle(options.SMARTctl.URLPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(options.SMARTctl.BindTo, nil))
}
