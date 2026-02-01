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
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfigFile(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{[]string{"--config.file", "config.yml"}, "config.yml"},
		{[]string{"--config.file=config.yml"}, "config.yml"},
		{[]string{"--other=flag"}, ""},
	}

	for _, test := range tests {
		if got := findConfigFile(test.args); got != test.expected {
			t.Fatalf("expected %q got %q", test.expected, got)
		}
	}
}

func TestLoadConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := []byte("web.listen-address: \":9999\"\nsmartctl.path: \"C:\\\\smartctl.exe\"\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := loadConfigFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.WebListenAddress == nil || *cfg.WebListenAddress != ":9999" {
		t.Fatalf("unexpected listen address: %#v", cfg.WebListenAddress)
	}
	if cfg.SmartctlPath == nil || *cfg.SmartctlPath != "C:\\smartctl.exe" {
		t.Fatalf("unexpected smartctl path: %#v", cfg.SmartctlPath)
	}
}
