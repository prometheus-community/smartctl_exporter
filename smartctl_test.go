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
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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

func TestMineDeviceSCTStatusReadsObjectValue(t *testing.T) {
	raw := `{
		"serial_number": "WD-TEST",
		"model_name": "WDC WD10EADS",
		"device": {"type": "sat", "protocol": "ATA"},
		"ata_sct_status": {
			"device_state": {"value": 4, "string": "SMART Off-line Data Collection executing in background"}
		}
	}`
	got := collectSCTStatus(t, raw)
	if got != 4 {
		t.Fatalf("smartctl_device_state = %v, want 4; Float() on the object would have been 0", got)
	}
}

func TestMineDeviceSCTStatusNumericFallback(t *testing.T) {
	raw := `{
		"serial_number": "WD-TEST",
		"model_name": "WDC WD10EADS",
		"device": {"type": "sat", "protocol": "ATA"},
		"ata_sct_status": {"device_state": 1}
	}`
	got := collectSCTStatus(t, raw)
	if got != 1 {
		t.Fatalf("smartctl_device_state = %v, want 1", got)
	}
}

func TestMineDeviceSCTStatusAbsent(t *testing.T) {
	raw := `{"serial_number":"WD-TEST","model_name":"WDC","device":{"type":"sat","protocol":"ATA"}}`
	ch := make(chan prometheus.Metric, 8)
	smart := NewSMARTctl(slog.Default(), gjson.Parse(raw), ch, Device{Name: "/dev/sda", Type: "sat", Label: "sda"})
	smart.mineDeviceSCTStatus()
	close(ch)
	for range ch {
		t.Fatal("did not expect a metric when ata_sct_status is missing")
	}
}

func collectSCTStatus(t *testing.T, raw string) float64 {
	t.Helper()
	ch := make(chan prometheus.Metric, 8)
	smart := NewSMARTctl(slog.Default(), gjson.Parse(raw), ch, Device{Name: "/dev/sda", Type: "sat", Label: "sda"})
	smart.mineDeviceSCTStatus()
	close(ch)
	var m prometheus.Metric
	for v := range ch {
		if m != nil {
			t.Fatalf("got more than one metric")
		}
		m = v
	}
	if m == nil {
		t.Fatal("expected a smartctl_device_state metric")
	}
	pb := &dto.Metric{}
	if err := m.Write(pb); err != nil {
		t.Fatalf("write metric: %v", err)
	}
	return pb.GetGauge().GetValue()
}
