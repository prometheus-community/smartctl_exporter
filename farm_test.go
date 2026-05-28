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
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

func TestMineFarmLog(t *testing.T) {
	jsonFile, err := os.ReadFile("testdata/sat-Seagate_Exos_X18-ST12000NM000J-farm.json")
	if err != nil {
		t.Fatal(err)
	}
	json := gjson.Parse(string(jsonFile))
	ch := make(chan prometheus.Metric, 1000)

	smart := NewSMARTctl(slog.Default(), json, ch)
	smart.mineFarmLog()
	close(ch)

	metrics := make(map[string][]prometheus.Metric)
	for m := range ch {
		metrics[m.Desc().String()] = append(metrics[m.Desc().String()], m)
	}

	// Verify workload metrics
	assertMetricCount(t, metrics, metricFarmWorkloadReadCommands, 1)
	assertMetricCount(t, metrics, metricFarmWorkloadWriteCommands, 1)
	assertMetricCount(t, metrics, metricFarmWorkloadLogicalSectorsRead, 1)
	assertMetricCount(t, metrics, metricFarmWorkloadLogicalSectorsWritten, 1)

	// Verify error metrics
	assertMetricCount(t, metrics, metricFarmErrorUnrecoverableRead, 1)
	assertMetricCount(t, metrics, metricFarmErrorUnrecoverableWrite, 1)
	assertMetricCount(t, metrics, metricFarmErrorReallocatedSectors, 1)
	assertMetricCount(t, metrics, metricFarmErrorReallocatedCandidates, 1)
	assertMetricCount(t, metrics, metricFarmErrorCRCErrors, 1)
	assertMetricCount(t, metrics, metricFarmErrorCommandTimeouts, 1)

	// Verify per-head error metrics: 16 heads * 2 types = 32
	assertMetricCount(t, metrics, metricFarmErrorUnrecoverableByHead, 32)

	// Verify environment metrics
	assertMetricCount(t, metrics, metricFarmEnvironmentTemperature, 9)
	assertMetricCount(t, metrics, metricFarmEnvironmentHumidity, 1)
	assertMetricCount(t, metrics, metricFarmEnvironmentMotorPower, 1)
	// 12V and 5V all have value 0 but still exist in JSON
	assertMetricCount(t, metrics, metricFarmEnvironment12V, 3)
	assertMetricCount(t, metrics, metricFarmEnvironment5V, 3)

	// Verify reliability metrics
	assertMetricCount(t, metrics, metricFarmReliabilityErrorRate, 1)
	assertMetricCount(t, metrics, metricFarmReliabilitySeekErrorRate, 1)
	assertMetricCount(t, metrics, metricFarmReliabilityHighPriorityUnloads, 1)
	assertMetricCount(t, metrics, metricFarmReliabilityHeliumPressureTrip, 1)

	// Per-head reliability: 16 heads each
	assertMetricCount(t, metrics, metricFarmReliabilityReallocatedByHead, 16)
	assertMetricCount(t, metrics, metricFarmReliabilityCandidatesByHead, 16)
	assertMetricCount(t, metrics, metricFarmReliabilityMRHeadResistanceByHead, 16)
	// Skip write: 16 heads * 4 types = 64
	assertMetricCount(t, metrics, metricFarmReliabilitySkipWriteByHead, 64)
}

func TestMineFarmLogAbsentKey(t *testing.T) {
	// Test that no FARM metrics are emitted when seagate_farm_log is missing
	json := gjson.Parse(`{"device": {"name": "/dev/sda", "type": "sat", "protocol": "ATA"}, "model_name": "WDC WD10EZEX"}`)
	ch := make(chan prometheus.Metric, 100)

	smart := NewSMARTctl(slog.Default(), json, ch)
	smart.mineFarmLog()
	close(ch)

	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 metrics for non-Seagate drive, got %d", count)
	}
}

func assertMetricCount(t *testing.T, metrics map[string][]prometheus.Metric, desc *prometheus.Desc, expected int) {
	t.Helper()
	actual := len(metrics[desc.String()])
	if actual != expected {
		t.Errorf("metric %s: expected %d, got %d", desc, expected, actual)
	}
}
