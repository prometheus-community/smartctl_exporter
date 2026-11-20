// Copyright 2022 The Prometheus Authors
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
	"path/filepath"
	"testing"
)

func TestParseExtraDeviceSpec(t *testing.T) {
	tests := []struct {
		input     string
		wantErr   bool
		wantDir   string
		wantRe    string
		wantType  string
	}{
		{"/dev/spdk;^nvme\\d+$;nvme", false, "/dev/spdk", `^nvme\d+$`, "nvme"},
		{"/dev/vfio;^[0-9]+$;nvme", false, "/dev/vfio", `^[0-9]+$`, "nvme"},
		{"missing-semicolons", true, "", "", ""},
		{"/dev/spdk;[invalid;nvme", true, "", "", ""},
		{"/dev/spdk;valid;", false, "/dev/spdk", "valid", ""},
	}
	for _, tt := range tests {
		spec, err := parseExtraDeviceSpec(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseExtraDeviceSpec(%q): expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseExtraDeviceSpec(%q): unexpected error: %v", tt.input, err)
			continue
		}
		if spec.directory != tt.wantDir {
			t.Errorf("parseExtraDeviceSpec(%q): directory = %q, want %q", tt.input, spec.directory, tt.wantDir)
		}
		if spec.nameRegexp.String() != tt.wantRe {
			t.Errorf("parseExtraDeviceSpec(%q): regexp = %q, want %q", tt.input, spec.nameRegexp.String(), tt.wantRe)
		}
		if spec.deviceType != tt.wantType {
			t.Errorf("parseExtraDeviceSpec(%q): deviceType = %q, want %q", tt.input, spec.deviceType, tt.wantType)
		}
	}
}

func TestScanExtraDevices(t *testing.T) {
	// Set up a temporary directory mimicking /dev/spdk with some nvme* entries.
	dir := t.TempDir()
	for _, name := range []string{"nvme0", "nvme1", "sda", "nvme2"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte{}, 0600); err != nil {
			t.Fatal(err)
		}
	}

	// Point global flags at empty strings so the filter accepts everything.
	empty := ""
	smartctlDeviceExclude = &empty
	smartctlDeviceInclude = &empty

	logger := slog.Default()
	specs := []string{dir + ";^nvme\\d+$;nvme"}
	devices := scanExtraDevices(logger, specs)

	if len(devices) != 3 {
		t.Fatalf("expected 3 nvme devices, got %d: %v", len(devices), devices)
	}
	for _, d := range devices {
		if d.Type != "nvme" {
			t.Errorf("device %q: expected type 'nvme', got %q", d.Name, d.Type)
		}
	}
}

func TestDeviceFilter(t *testing.T) {
	tests := []struct {
		ignore         string
		accept         string
		name           string
		expectedResult bool
	}{
		{"", "", "eth0", false},
		{"", "^💩0$", "💩0", false},
		{"", "^💩0$", "💩1", true},
		{"", "^💩0$", "veth0", true},
		{"^💩", "", "💩3", true},
		{"^💩", "", "veth0", false},
	}

	for _, test := range tests {
		filter := newDeviceFilter(test.ignore, test.accept)
		result := filter.ignored(test.name)

		if result != test.expectedResult {
			t.Errorf("ignorePattern=%v acceptPattern=%v ifname=%v expected=%v result=%v", test.ignore, test.accept, test.name, test.expectedResult, result)
		}
	}
}
