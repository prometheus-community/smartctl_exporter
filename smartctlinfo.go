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
	smart.mineVersion()
}

func (smart *SMARTctlInfo) mineVersion() {
	smartctlJSON := smart.json.Get("smartctl")
	smartctlVersion := smartctlJSON.Get("version").Array()
	jsonVersion := smart.json.Get("json_format_version").Array()
	if len(smartctlVersion) > 0 && len(jsonVersion) > 0 {
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
}
