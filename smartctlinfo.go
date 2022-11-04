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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

// SMARTctlInfo object
type SMARTctlInfo struct {
	ch    chan<- prometheus.Metric
	json  gjson.Result
	Ready bool
}

// NewSMARTctlInfo is smartctl constructor
func NewSMARTctlInfo(ch chan<- prometheus.Metric) SMARTctlInfo {
	smart := SMARTctlInfo{}
	smart.ch = ch
	smart.Ready = false
	return smart
}

// SetJSON metrics
func (smart *SMARTctlInfo) SetJSON(json gjson.Result) {
	if !smart.Ready {
		smart.json = json
		smart.Ready = true
	}
}

// Collect metrics
func (smart *SMARTctlInfo) Collect() {
	if smart.Ready {
		smart.mineVersion()
	}
}

func (smart *SMARTctlInfo) mineVersion() {
	smartctlJSON := smart.json.Get("smartctl")
	smartctlVersion := smartctlJSON.Get("version").Array()
	jsonVersion := smart.json.Get("json_format_version").Array()
	smart.ch <- prometheus.MustNewConstMetric(
		metricSmartctlVersion,
		prometheus.GaugeValue,
		1,
		fmt.Sprintf("%d.%d", jsonVersion[0].Int(), jsonVersion[1].Int()),
		fmt.Sprintf("%d.%d", smartctlVersion[0].Int(), smartctlVersion[1].Int()),
		smartctlJSON.Get("svn_revision").String(),
		smartctlJSON.Get("build_info").String(),
	)
}
