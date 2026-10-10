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
	"fmt"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/tidwall/gjson"
)

func TestCollectSCSIMetrics(t *testing.T) {
	want := map[*prometheus.Desc]float64{
		metricSCSIGrownDefectList:                    20,
		metricSCSIUsedEnduranceIndicator:             7,
		metricReadErrorsCorrectedByRereadsRewrites:   1,
		metricReadErrorsCorrectedByEccFast:           2,
		metricReadErrorsCorrectedByEccDelayed:        3,
		metricReadTotalUncorrectedErrors:             4,
		metricWriteErrorsCorrectedByRereadsRewrites:  5,
		metricWriteErrorsCorrectedByEccFast:          6,
		metricWriteErrorsCorrectedByEccDelayed:       7,
		metricWriteTotalUncorrectedErrors:            8,
		metricVerifyErrorsCorrectedByRereadsRewrites: 9,
		metricVerifyErrorsCorrectedByEccFast:         10,
		metricVerifyErrorsCorrectedByEccDelayed:      11,
		metricVerifyTotalUncorrectedErrors:           12,
		metricDeviceBytesRead:                        42.5e9,
		metricDeviceBytesWritten:                     25.5e9,
	}
	for _, test := range []struct {
		name       string
		deviceType string
		protocol   string
		wantSCSI   bool
	}{
		{"native SCSI", "scsi", "SCSI", true},
		{"SCSI without protocol", "scsi", "", true},
		{"MegaRAID SCSI", "megaraid,0", "SCSI", true},
		{"MegaRAID ATA", "megaraid,0", "ATA", false},
		{"SAT MegaRAID ATA", "sat+megaraid,0", "ATA", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Keep the SCSI fields in the ATA cases to verify the protocol gate.
			data := gjson.Parse(fmt.Sprintf(`{
				"smartctl": {"exit_status": 0},
				"device": {"type": %q, "protocol": %q},
				"scsi_grown_defect_list": 20,
				"scsi_percentage_used_endurance_indicator": 7,
				"scsi_error_counter_log": {
					"read": {
						"errors_corrected_by_rereads_rewrites": 1,
						"errors_corrected_by_eccfast": 2,
						"errors_corrected_by_eccdelayed": 3,
						"total_uncorrected_errors": 4,
						"gigabytes_processed": "42.5"
					},
					"write": {
						"errors_corrected_by_rereads_rewrites": 5,
						"errors_corrected_by_eccfast": 6,
						"errors_corrected_by_eccdelayed": 7,
						"total_uncorrected_errors": 8,
						"gigabytes_processed": "25.5"
					},
					"verify": {
						"errors_corrected_by_rereads_rewrites": 9,
						"errors_corrected_by_eccfast": 10,
						"errors_corrected_by_eccdelayed": 11,
						"total_uncorrected_errors": 12
					}
				}
			}`, test.deviceType, test.protocol))
			ch := make(chan prometheus.Metric)
			logger := slog.New(slog.DiscardHandler)
			smart := NewSMARTctl(logger, data, ch, Device{Label: "test-disk"})
			go func() {
				smart.Collect()
				close(ch)
			}()

			// Inspect device collection directly, independently of the manager's
			// descriptor registration.
			seen := make(map[*prometheus.Desc]bool)
			for metric := range ch {
				desc := metric.Desc()
				value, ok := want[desc]
				if !ok {
					continue
				}
				if seen[desc] {
					t.Errorf("duplicate metric %s", desc)
				}
				seen[desc] = true
				var sample dto.Metric
				if err := metric.Write(&sample); err != nil {
					t.Errorf("write metric %s: %v", desc, err)
					continue
				}
				got := sample.GetGauge().GetValue()
				if desc == metricDeviceBytesRead || desc == metricDeviceBytesWritten {
					got = sample.GetCounter().GetValue()
				}
				if got != value {
					t.Errorf("%s: got %v, want %v", desc, got, value)
				}
				labels := sample.GetLabel()
				if len(labels) != 1 || labels[0].GetName() != "device" || labels[0].GetValue() != "test-disk" {
					t.Errorf("%s: unexpected labels %v", desc, labels)
				}
			}
			for desc := range want {
				if seen[desc] != test.wantSCSI {
					t.Errorf("%s: collected=%v, want %v", desc, seen[desc], test.wantSCSI)
				}
			}
		})
	}
}
