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
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
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

func TestStandbyJSONSkipsMissingMetrics(t *testing.T) {
	names := collectMetricNames(t, "testdata/standby-sdc.json")
	expected := map[string]struct{}{
		"smartctl_device_power_mode":           {},
		"smartctl_device_smartctl_exit_status": {},
	}

	if len(names) != len(expected) {
		t.Fatalf("expected %d metrics, got %d: %v", len(expected), len(names), names)
	}
	for name := range expected {
		if _, ok := names[name]; !ok {
			t.Fatalf("missing metric %q", name)
		}
	}
	for name := range names {
		if _, ok := expected[name]; !ok {
			t.Fatalf("unexpected metric %q", name)
		}
	}
}

func collectMetricNames(t *testing.T, jsonPath string) map[string]struct{} {
	t.Helper()

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("read json: %v", err)
	}
	json := gjson.ParseBytes(data)

	ch := make(chan prometheus.Metric)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	smart := NewSMARTctl(logger, json, ch)

	go func() {
		smart.Collect()
		close(ch)
	}()

	names := make(map[string]struct{})
	for metric := range ch {
		name := metricName(metric)
		if name == "" {
			t.Fatalf("missing metric name for %v", metric)
		}
		names[name] = struct{}{}
	}
	return names
}

func metricName(metric prometheus.Metric) string {
	desc := metric.Desc().String()
	const prefix = `fqName: "`
	start := strings.Index(desc, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	end := strings.Index(desc[start:], `"`)
	if end == -1 {
		return ""
	}
	return desc[start : start+end]
}
