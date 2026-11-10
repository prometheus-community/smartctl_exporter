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
	"path/filepath"
	"sync"
	"testing"
	"time"
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

// fakeSmartctl installs a shell script as smartctl and applies timeout.
func fakeSmartctl(t *testing.T, body string, timeout time.Duration) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "smartctl")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	oldPath, oldTimeout := *smartctlPath, *smartctlTimeout
	*smartctlPath, *smartctlTimeout = path, timeout
	t.Cleanup(func() { *smartctlPath, *smartctlTimeout = oldPath, oldTimeout })
}

func readOnce(device Device) {
	var wg sync.WaitGroup
	wg.Add(1)
	readSMARTctl(slog.New(slog.DiscardHandler), device, &wg)
}

func scrapeOk(device Device) (ok, found bool) {
	v, found := deviceScrapeOk.Load(device)
	if !found {
		return false, false
	}
	return v.(bool), true
}

func TestReadSMARTctlTimeout(t *testing.T) {
	device := Device{Name: "/dev/hung", Type: "auto", Label: "hung"}
	good := `echo '{"smartctl":{"exit_status":0}}'`

	fakeSmartctl(t, good, time.Second)
	readOnce(device)
	if _, cached := jsonCache.Load(device); !cached {
		t.Fatal("expected data to be cached after successful read")
	}

	// Hung smartctl: must return after the timeout and drop the stale data.
	// Not covered: a process in D state that ignores SIGKILL.
	fakeSmartctl(t, "exec sleep 30", 300*time.Millisecond)
	start := time.Now()
	readOnce(device)
	if d := time.Since(start); d > 5*time.Second {
		t.Fatalf("readSMARTctl took %v, timeout not applied", d)
	}
	if _, cached := jsonCache.Load(device); cached {
		t.Error("stale data must be dropped after a failed read")
	}
	if ok, found := scrapeOk(device); !found || ok {
		t.Errorf("expected scrape failure, got ok=%v found=%v", ok, found)
	}

	// Killed process eventually exits and releases the in-flight flag.
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, running := inFlight.Load(device); !running {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("in-flight flag not released after process exit")
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Recovery repopulates the cache.
	fakeSmartctl(t, good, time.Second)
	readOnce(device)
	if _, cached := jsonCache.Load(device); !cached {
		t.Error("expected cache to be repopulated after recovery")
	}
	if ok, _ := scrapeOk(device); !ok {
		t.Error("expected scrape success after recovery")
	}
}

func TestReadSMARTctlStandbyKeepsCache(t *testing.T) {
	device := Device{Name: "/dev/standby", Type: "auto", Label: "standby"}
	fakeSmartctl(t, `echo '{"smartctl":{"exit_status":0}}'`, time.Second)
	readOnce(device)

	fakeSmartctl(t, `echo '{"smartctl":{"exit_status":2,"messages":[{"severity":"information","string":"standby"}]}}'; exit 2`, time.Second)
	readOnce(device)
	if _, cached := jsonCache.Load(device); !cached {
		t.Error("standby must not drop cached data")
	}
}
