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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	// "github.com/golang/protobuf/proto"
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

func TestMineDeviceSelfTestLog(t *testing.T) {
	tests := []struct {
		name     string
		jsonFile string
		want     map[string]struct {
			count      float64
			errorTotal float64
			logType    string
		}
	}{
		{
			name:     "Exos X16 self-test log parsing",
			jsonFile: "testdata/sat-Segate_Exos_X16-ST10000NM001G-2MW103.json",
			want: map[string]struct {
				count      float64
				errorTotal float64
				logType    string
			}{
				"standard": {
					count:      21,
					errorTotal: 0,
					logType:    "standard",
				},
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

			// Verify metrics
			metricsFound := make(map[string]bool)
			for m := range ch {
				metric := new(dto.Metric)
				m.Write(metric)

				var deviceName, logType string
				for _, label := range metric.GetLabel() {
					switch *label.Name {
					case "device":
						deviceName = *label.Value
					case "self_test_log_type":
						logType = *label.Value
					}
				}

				// Verify device name matches JSON contents
				if deviceName != "sdc" {
					t.Errorf("Unexpected device name: got %s want sdc", deviceName)
				}

				// Check log type and values
				expected, ok := tt.want[logType]
				if !ok {
					t.Errorf("Unexpected log type: %s", logType)
					continue
				}

				metricType := getMetricType(metric)
				if metricType == nil {
					t.Errorf("Unknown metric type")
				}
				switch *metricType {
				case dto.MetricType_GAUGE:
					value := metric.GetGauge().GetValue()
					if logType == "standard" {
						switch m.Desc() {
						case metricDeviceSelfTestLogCount:
							if value != expected.count {
								t.Errorf("count mismatch: got %v want %v", value, expected.count)
							}
						case metricDeviceSelfTestLogErrorCount:
							if value != expected.errorTotal {
								t.Errorf("error_count_total mismatch: got %v want %v", value, expected.errorTotal)
							}
						default:
							t.Errorf("Metric unkown")
						}

					}
					metricsFound[logType] = true
				}
			}

			// Verify we found all expected log types
			for logType := range tt.want {
				if !metricsFound[logType] {
					t.Errorf("Missing metrics for log type: %s", logType)
				}
			}
		})
	}
}

// Helper to determine metric type
func getMetricType(m *dto.Metric) *dto.MetricType {
	if m.Counter != nil {
		return dto.MetricType_COUNTER.Enum()
	}
	if m.Gauge != nil {
		return dto.MetricType_GAUGE.Enum()
	}
	if m.Histogram != nil {
		return dto.MetricType_HISTOGRAM.Enum()
	}
	if m.Summary != nil {
		return dto.MetricType_SUMMARY.Enum()
	}
	return nil
}
