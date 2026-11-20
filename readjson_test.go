// Copyright The Prometheus Authors
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
	"encoding/json"
	"log/slog"
	"os"
	"slices"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// Run the test binary as a fake smartctl, avoiding hardware dependencies and
// platform-specific shell scripts. SSD statistics require --log=devstat,7.
func TestMain(m *testing.M) {
	if os.Getenv("SMARTCTL_EXPORTER_TEST_SMARTCTL") == "1" {
		data := map[string]any{
			"json_format_version": []int{1, 0},
			"smartctl": map[string]any{
				"version":     []int{7, 5},
				"exit_status": 0,
				"argv":        os.Args[1:],
			},
			"device":       map[string]string{"name": "/dev/test SSD", "type": "atacam", "protocol": "ATA"},
			"smart_status": map[string]bool{"passed": true},
			"ata_smart_error_log": map[string]any{
				"summary": map[string]int{"count": 3},
			},
		}
		if slices.Contains(os.Args[1:], "--log=devstat,7") && os.Getenv("SMARTCTL_EXPORTER_TEST_NO_SSD_STATISTICS") != "1" {
			data["ata_device_statistics"] = map[string]any{
				"pages": []any{map[string]any{
					"name": "Solid State Device Statistics",
					"table": []any{map[string]any{
						"name":  "Percentage Used Endurance Indicator",
						"value": 10,
						"flags": map[string]any{"string": "V--- ", "valid": true},
					}},
				}},
			}
		}
		if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
			os.Exit(1)
		}
		return
	}
	os.Exit(m.Run())
}

func useFakeSMARTctl(t *testing.T) *slog.Logger {
	t.Helper()
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("SMARTCTL_EXPORTER_TEST_SMARTCTL", "1")
	oldPath, oldSSDStatistics, oldPowerMode := *smartctlPath, *smartctlSSDStatistics, *smartctlPowerModeCheck
	t.Cleanup(func() {
		*smartctlPath, *smartctlSSDStatistics, *smartctlPowerModeCheck = oldPath, oldSSDStatistics, oldPowerMode
	})
	*smartctlPath = helper
	return slog.New(slog.DiscardHandler)
}

func TestReadSMARTctlSSDStatistics(t *testing.T) {
	logger := useFakeSMARTctl(t)

	for _, test := range []struct {
		name          string
		ssdStatistics bool
		powerMode     string
		unsupported   bool
	}{
		{name: "default", powerMode: "standby"},
		{name: "ssd_statistics", ssdStatistics: true, powerMode: "standby"},
		{name: "custom_power_mode", ssdStatistics: true, powerMode: "idle"},
		{name: "unsupported_statistics", ssdStatistics: true, powerMode: "standby", unsupported: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			*smartctlSSDStatistics, *smartctlPowerModeCheck = test.ssdStatistics, test.powerMode
			if test.unsupported {
				t.Setenv("SMARTCTL_EXPORTER_TEST_NO_SSD_STATISTICS", "1")
			}
			device := Device{Name: "/dev/test SSD", Type: "atacam", Label: "test_SSD"}
			jsonCache.Delete(device)
			t.Cleanup(func() { jsonCache.Delete(device) })
			var wg sync.WaitGroup
			wg.Add(1)
			readSMARTctl(logger, device, &wg)
			wg.Wait()

			cached, ok := jsonCache.Load(device)
			if !ok {
				t.Fatal("smartctl output was not cached")
			}
			var args []string
			for _, arg := range cached.(JSONCache).JSON.Get("smartctl.argv").Array() {
				args = append(args, arg.String())
			}
			wantArgs := []string{"--json", "--info", "--health", "--attributes", "--capabilities", "--tolerance=verypermissive", "--nocheck=" + test.powerMode, "--format=brief", "--log=error"}
			if test.ssdStatistics {
				wantArgs = append(wantArgs, "--log=devstat,7")
			}
			wantArgs = append(wantArgs, "--device=atacam", "/dev/test SSD")
			if !slices.Equal(args, wantArgs) {
				t.Fatalf("smartctl arguments = %q, want %q", args, wantArgs)
			}

			collector := &SMARTctlManagerCollector{Devices: []Device{device}, logger: logger}
			registry := prometheus.NewPedanticRegistry()
			registry.MustRegister(collector)
			families, err := registry.Gather()
			if err != nil {
				t.Fatal(err)
			}
			values := make(map[string]float64)
			for _, family := range families {
				for _, metric := range family.GetMetric() {
					if metric.GetGauge() != nil {
						values[family.GetName()] = metric.GetGauge().GetValue()
					}
					if family.GetName() == "smartctl_device_statistics" {
						labels := make(map[string]string)
						for _, label := range metric.GetLabel() {
							labels[label.GetName()] = label.GetValue()
						}
						if labels["device"] != device.Label || labels["statistic_name"] != "Percentage Used Endurance Indicator" || labels["statistic_flags_long"] != "valid" {
							t.Errorf("unexpected device statistics labels: %v", labels)
						}
					}
				}
			}
			for name, want := range map[string]float64{
				"smartctl_device_smart_status":    1,
				"smartctl_device_error_log_count": 3,
			} {
				if got, exists := values[name]; !exists || got != want {
					t.Errorf("%s = %v (present=%v), want %v", name, got, exists, want)
				}
			}
			got, exists := values["smartctl_device_statistics"]
			wantStatistics := test.ssdStatistics && !test.unsupported
			if exists != wantStatistics || (exists && got != 10) {
				t.Errorf("endurance used = %v (present=%v), want 10 (present=%v)", got, exists, wantStatistics)
			}
		})
	}
}

func TestSSDStatisticsDoesNotChangeDeviceScan(t *testing.T) {
	logger := useFakeSMARTctl(t)
	*smartctlSSDStatistics = true
	oldTypes := *smartctlScanDeviceTypes
	t.Cleanup(func() { *smartctlScanDeviceTypes = oldTypes })
	*smartctlScanDeviceTypes = []string{"sat", "nvme"}

	result := readSMARTctlDevices(logger)
	var args []string
	for _, arg := range result.Get("smartctl.argv").Array() {
		args = append(args, arg.String())
	}
	want := []string{"--json", "--scan", "--device", "sat", "--device", "nvme"}
	if !slices.Equal(args, want) {
		t.Fatalf("scan arguments = %q, want %q", args, want)
	}
}
