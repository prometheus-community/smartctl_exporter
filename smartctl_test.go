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
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

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

func getLabelValue(labels []*dto.LabelPair, key string) (string, bool) {
	for _, label := range labels {
		if label.GetName() == key {
			return label.GetValue(), true
		}
	}
	return "", false
}

func getMetricsFromChannel(ch chan prometheus.Metric) map[*prometheus.Desc]*dto.Metric {
	metricMap := make(map[*prometheus.Desc]*dto.Metric)
	for m := range ch {
		metric := new(dto.Metric)
		m.Write(metric)
		metricMap[m.Desc()] = metric
	}
	return metricMap
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read and parse test JSON
			jsonRaw := readTestFile(tt.jsonFile)
			jsonData := parseJSON(string(jsonRaw))

			// Create collector and mine data
			ch := make(chan prometheus.Metric, 20) // Increased buffer size
			smart := NewSMARTctl(nil, jsonData, ch)
			smart.mineDeviceSelfTestLog()
			close(ch)

			metricMap := getMetricsFromChannel(ch)
			expected := tt.want

			metric := metricMap[metricDeviceSelfTestLogCount]
			assert.NotNil(t, metric, "Missing metricDeviceSelfTestLogCount")
			assert.Equal(t, expected.count, metric.GetGauge().GetValue())
			val, ok := getLabelValue(metric.GetLabel(), "device")
			assert.True(t, ok)
			assert.Equal(t, "sdc", val)
			val, ok = getLabelValue(metric.GetLabel(), "self_test_log_type")
			assert.True(t, ok)
			assert.Equal(t, "standard", val)

			metric = metricMap[metricDeviceSelfTestLogErrorCount]
			assert.NotNil(t, metric, "Missing metricDeviceSelfTestLogErrorCount")
			assert.Equal(t, expected.errorTotal, metric.GetGauge().GetValue())
			val, ok = getLabelValue(metric.GetLabel(), "device")
			assert.True(t, ok)
			assert.Equal(t, "sdc", val)
			val, ok = getLabelValue(metric.GetLabel(), "self_test_log_type")
			assert.True(t, ok)
			assert.Equal(t, "standard", val)

			metric = metricMap[metricDeviceLastSelfTest]
			assert.NotNil(t, metric, "Missing metricDeviceLastSelfTest")
			assert.Equal(t, expected.lastTestStatus, metric.GetGauge().GetValue())
			val, ok = getLabelValue(metric.GetLabel(), "device")
			assert.True(t, ok)
			assert.Equal(t, "sdc", val)
			val, ok = getLabelValue(metric.GetLabel(), "lifetime_hours")
			assert.True(t, ok)
			assert.Equal(t, expected.lastTestHours, val)

			metric = metricMap[metricDeviceLastSelfTestInfo]
			assert.NotNil(t, metric, "Missing metricDeviceLastSelfTestInfo")
			assert.Equal(t, 1.0, metric.GetGauge().GetValue())
			val, ok = getLabelValue(metric.GetLabel(), "device")
			assert.True(t, ok)
			assert.Equal(t, "sdc", val)
			val, ok = getLabelValue(metric.GetLabel(), "lifetime_hours")
			assert.True(t, ok)
			assert.Equal(t, expected.lastTestHours, val)
			val, ok = getLabelValue(metric.GetLabel(), "status")
			assert.True(t, ok)
			assert.Equal(t, strconv.FormatFloat(expected.lastTestStatus, 'f', -1, 64), val)

			val, ok = getLabelValue(metric.GetLabel(), "description")
			assert.True(t, ok)
			assert.Equal(t, expected.lastTestStatusDesc, val)
			val, ok = getLabelValue(metric.GetLabel(), "type")
			assert.True(t, ok)
			assert.Equal(t, expected.lastTestType, val)

		})
	}
}
