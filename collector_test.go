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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promslog"
)

// fixtureCollector describes itself like SMARTctlManagerCollector and collects
// the metrics mined from smartctl JSON fixtures.
type fixtureCollector struct {
	manager  *SMARTctlManagerCollector
	fixtures map[string]string
}

func (c *fixtureCollector) Describe(ch chan<- *prometheus.Desc) {
	c.manager.Describe(ch)
}

func (c *fixtureCollector) Collect(ch chan<- prometheus.Metric) {
	for name, data := range c.fixtures {
		device := Device{Name: "/dev/" + name, Type: "auto", Label: name}
		smart := NewSMARTctl(c.manager.logger, parseJSON(data), ch, device)
		smart.Collect()
	}
}

func loadFixtures(t *testing.T) map[string]string {
	t.Helper()
	paths, err := filepath.Glob("testdata/*.json")
	if err != nil || len(paths) == 0 {
		t.Fatalf("no fixtures found: %v", err)
	}
	fixtures := make(map[string]string, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		fixtures[strings.TrimSuffix(filepath.Base(path), ".json")] = string(data)
	}
	return fixtures
}

// TestCollectedMetricsAreDescribed fails when a mined metric is missing from
// Describe: the registry then rejects the whole scrape with HTTP 500.
func TestCollectedMetricsAreDescribed(t *testing.T) {
	fixtures := loadFixtures(t)
	// No fixture reports the SCSI endurance indicator; add it to a SAS one.
	sas := fixtures["HITACHI_H109060SESUN600G_10"]
	if sas == "" {
		t.Fatal("SAS fixture HITACHI_H109060SESUN600G_10 missing")
	}
	fixtures["sas-endurance"] = strings.Replace(sas, "{", `{"scsi_percentage_used_endurance_indicator": 3,`, 1)

	manager := &SMARTctlManagerCollector{logger: promslog.NewNopLogger()}
	registry := prometheus.NewPedanticRegistry()
	if err := registry.Register(&fixtureCollector{manager: manager, fixtures: fixtures}); err != nil {
		t.Fatalf("register: %v", err)
	}

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	gathered := make(map[string]bool, len(families))
	for _, family := range families {
		gathered[family.GetName()] = true
	}
	for _, name := range []string{
		"smartctl_scsi_percentage_used_endurance_indicator",
		"smartctl_verify_total_uncorrected_errors",
	} {
		if !gathered[name] {
			t.Errorf("metric %s not gathered", name)
		}
	}
}
