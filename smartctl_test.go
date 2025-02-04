// Copyright 2024 The Prometheus Authors
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
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

// channelCollector adapts a channel-based metric producer to Prometheus Collector interface
type channelCollector struct {
	metrics <-chan prometheus.Metric
}

func (c *channelCollector) Describe(ch chan<- *prometheus.Desc) {
	// No-op implementation since we're dealing with dynamic metrics from smartctl.
}

func (c *channelCollector) Collect(ch chan<- prometheus.Metric) {
	// Forward metrics from our internal channel to Prometheus
	for m := range c.metrics {
		ch <- m
	}
}

type MetricFamilies []*io_prometheus_client.MetricFamily

// gathermetrics handles metric collection for testing purposes using Prometheus registry
func gathermetrics(t *testing.T, ch <-chan prometheus.Metric) MetricFamilies {
	reg := prometheus.NewRegistry()

	collector := &channelCollector{metrics: ch}
	if err := reg.Register(collector); err != nil {
		t.Fatalf("Failed to register collector: %v", err)
	}

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Encode gathered metrics into the text format.
	var buf bytes.Buffer
	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, mf := range mfs {
		if err := enc.Encode(mf); err != nil {
			t.Fatalf("failed to encode metric family: %v", err)
		}
	}
	output := buf.String()
	t.Log("Gathered metrics output:\n", output)

	return mfs
}

func (m MetricFamilies) GetMetricFamily(name string) (*io_prometheus_client.MetricFamily, error) {
	for _, mf := range m {
		if mf.GetName() == name {
			return mf, nil
		}
	}
	return nil, fmt.Errorf("metric family %s not found", name)
}

func (m MetricFamilies) GetMetricWithLabelMap(name string, labels map[string]string) (*io_prometheus_client.Metric, error) {
	// Look up the metric family by name.
	family, ok := m.GetMetricFamily(name)
	if ok != nil {
		return nil, fmt.Errorf("metric family %q not found", name)
	}

	// Iterate over each metric in the family.
	for _, metric := range family.Metric {
		matches := true

		// For each key/value pair in the labels map, check if it exists in the metric.
		for key, val := range labels {
			labelFound := false

			// Check the metric's labels.
			for _, lp := range metric.Label {
				if lp.GetName() == key {
					if lp.GetValue() == val {
						labelFound = true
					}
					break // Found the label key, break out of the inner loop.
				}
			}

			// If the current label was not found or didn't match the expected value, skip this metric.
			if !labelFound {
				matches = false
				break
			}
		}

		// Return the metric if all provided labels match.
		if matches {
			return metric, nil
		}
	}

	return nil, fmt.Errorf("metric %q with labels %v not found", name, labels)
}

func TestBuildDeviceLabel(t *testing.T) {
	tests := []struct {
		deviceName    string
		deviceType    string
		expectedLabel string
	}{
		{"/dev/bus/0", "megaraid,1", "bus_0_megaraid_1"},
		{"/dev/sda", "auto", "sda"},
		{"/dev/disk/by-id/ata-CT500MX500SSD1_ABCDEFGHIJ", "auto", "ata-CT500MX500SSD1_ABCDEFGHIJ"},
		// Some cases extracted from smartctl docs. Are these the prettiest?
		// Probably not. Are they unique enough. Definitely.
		{"/dev/sg1", "cciss,1", "sg1_cciss_1"},
		{"/dev/bsg/sssraid0", "sssraid,0,1", "bsg_sssraid0_sssraid_0_1"},
		{"/dev/cciss/c0d0", "cciss,0", "cciss_c0d0_cciss_0"},
		{"/dev/sdb", "aacraid,1,0,4", "sdb_aacraid_1_0_4"},
		{"/dev/twl0", "3ware,1", "twl0_3ware_1"},
	}

	for _, test := range tests {
		result := buildDeviceLabel(test.deviceName, test.deviceType)
		if result != test.expectedLabel {
			t.Errorf("deviceName=%v deviceType=%v expected=%v result=%v", test.deviceName, test.deviceType, test.expectedLabel, result)
		}
	}
}

func readTestFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading test file: %v", err))
	}
	return data
}

func TestMineDeviceSelfTestLog(t *testing.T) {
	tests := []struct {
		name     string
		jsonFile string
		want     struct {
			count              float64
			errorTotal         float64
			logType            string
			lastTestType       string
			lastTestHours      string
			lastTestStatus     float64
			lastTestStatusDesc string
		}
	}{
		{
			name:     "Exos X16 self-test log parsing",
			jsonFile: "testdata/sat-Segate_Exos_X16-ST10000NM001G-2MW103.json",
			want: struct {
				count              float64
				errorTotal         float64
				logType            string
				lastTestType       string
				lastTestHours      string
				lastTestStatus     float64
				lastTestStatusDesc string
			}{

				count:              21,
				errorTotal:         0,
				logType:            "standard",
				lastTestType:       "Short offline",
				lastTestHours:      "33600",
				lastTestStatus:     0,
				lastTestStatusDesc: "Completed without error",
			},
		},
		{
			name:     "Scsi-seagate_ST18000NM004J",
			jsonFile: "testdata/scsi-seagate_ST18000NM004J.json",
			want: struct {
				count              float64
				errorTotal         float64
				logType            string
				lastTestType       string
				lastTestHours      string
				lastTestStatus     float64
				lastTestStatusDesc string
			}{

				count:              -999,
				errorTotal:         -999,
				logType:            "",
				lastTestType:       "Background short",
				lastTestHours:      "1239",
				lastTestStatus:     0,
				lastTestStatusDesc: "Completed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read and parse test JSON
			jsonRaw := readTestFile(tt.jsonFile)
			jsonData := parseJSON(string(jsonRaw))

			// Extract device name from JSON
			deviceName := jsonData.Get("device.name").String()
			deviceName = strings.TrimPrefix(deviceName, "/dev/")

			// Create collector and mine data
			ch := make(chan prometheus.Metric, 20)
			smart := NewSMARTctl(nil, jsonData, ch)
			smart.mineDeviceSelfTestLog()
			close(ch)

			// Get registry and metrics
			mfs := gathermetrics(t, ch)

			expected := tt.want

			// Validate self test count metric if expected
			if expected.count >= 0 {
				metric, err := mfs.GetMetricWithLabelMap("smartctl_device_self_test_log_count",
					map[string]string{"device": deviceName, "self_test_log_type": expected.logType})
				assert.NoError(t, err, "metric smartctl_device_self_test_log_count not found")
				assert.Equal(t, expected.count, metric.GetGauge().GetValue(), "metric smartctl_device_self_test_log_count value")
			}

			// Execute if we expect an error total metric
			if expected.errorTotal >= 0 {
				metric, err := mfs.GetMetricWithLabelMap("smartctl_device_self_test_log_error_count",
					map[string]string{"device": deviceName, "self_test_log_type": expected.logType})
				assert.NoError(t, err, "metric smartctl_device_self_test_log_error_count not found")
				assert.Equal(t, expected.errorTotal, metric.GetGauge().GetValue(), "metric smartctl_device_self_test_log_error_count value")
			}

			metric, err := mfs.GetMetricWithLabelMap("smartctl_device_last_self_test",
				map[string]string{"device": deviceName, "lifetime_hours": expected.lastTestHours})
			assert.NoError(t, err, "metric smartctl_device_last_self_test not found")
			assert.Equal(t, expected.lastTestStatus, metric.GetGauge().GetValue(), "metric smartctl_device_last_self_test value")

			metric, err = mfs.GetMetricWithLabelMap("smartctl_device_last_self_test_info",
				map[string]string{"device": deviceName, "lifetime_hours": expected.lastTestHours, "status": strconv.FormatFloat(expected.lastTestStatus, 'f', -1, 64), "description": expected.lastTestStatusDesc})
			assert.NoError(t, err, "metric smartctl_device_last_self_test_info not found")
			assert.Equal(t, 1.0, metric.GetGauge().GetValue(), "metric smartctl_device_last_self_test_info value")
		})
	}
}
