// Copyright 2026 The Prometheus Authors
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
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

// smartctl writes each of read/write/verify in scsi_error_counter_log only
// when the drive supports that error counter log page.
func TestSCSIErrorCounterLogMissingVerify(t *testing.T) {
	// Redacted smartctl 7.5 output for a SAS drive without the verify error
	// counter log page, from
	// https://github.com/prometheus-community/smartctl_exporter/issues/406
	data, err := os.ReadFile("testdata/scsi-no-verify-counters.json")
	if err != nil {
		t.Fatal(err)
	}
	json := gjson.ParseBytes(data)

	ch := make(chan prometheus.Metric)
	smart := NewSMARTctl(slog.New(slog.DiscardHandler), json, ch, Device{Name: "/dev/sdbv", Type: "scsi", Label: "sdbv"})
	go func() {
		smart.Collect()
		close(ch)
	}()
	got := map[*prometheus.Desc]bool{}
	for m := range ch {
		got[m.Desc()] = true
	}

	for _, desc := range []*prometheus.Desc{
		metricDeviceBytesRead,
		metricReadErrorsCorrectedByRereadsRewrites,
		metricReadErrorsCorrectedByEccFast,
		metricReadErrorsCorrectedByEccDelayed,
		metricReadTotalUncorrectedErrors,
		metricDeviceBytesWritten,
		metricWriteErrorsCorrectedByRereadsRewrites,
		metricWriteErrorsCorrectedByEccFast,
		metricWriteErrorsCorrectedByEccDelayed,
		metricWriteTotalUncorrectedErrors,
	} {
		if !got[desc] {
			t.Errorf("missing %s", desc)
		}
	}
	for _, desc := range []*prometheus.Desc{
		metricVerifyErrorsCorrectedByRereadsRewrites,
		metricVerifyErrorsCorrectedByEccFast,
		metricVerifyErrorsCorrectedByEccDelayed,
		metricVerifyTotalUncorrectedErrors,
	} {
		if got[desc] {
			t.Errorf("unexpected %s", desc)
		}
	}
}
