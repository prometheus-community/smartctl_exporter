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
	"testing"

	"github.com/tidwall/gjson"
)

func TestResultCodeIsOkStandbyJSON(t *testing.T) {
	// output from a standby hard drive:
	// sudo hdparm -y /dev/sdc
	// sudo smartctl --nocheck=standby /dev/sdc --json --info --health --attributes --tolerance=verypermissive --format=brief --log=error
	json := readTestJSON(t, "testdata/standby-sdc.json")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	device := Device{
		Name:  "/dev/sdc",
		Type:  "sat",
		Label: "sdc",
	}

	if !resultCodeIsOk(logger, device, json.Get("smartctl.exit_status").Int(), json) {
		t.Fatalf("expected exit status to be ok for standby json")
	}
}

func TestResultCodeIsOkNonexistentDeviceJSON(t *testing.T) {
	// output from a nonexistent disk:
	// sudo smartctl --nocheck=standby /dev/nonexistent --json --info --health --attributes --tolerance=verypermissive --format=brief --log=error
	json := readTestJSON(t, "testdata/nonexistent-device.json")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	device := Device{
		Name:  "/dev/nonexistent",
		Type:  "auto",
		Label: "nonexistent",
	}

	if resultCodeIsOk(logger, device, json.Get("smartctl.exit_status").Int(), json) {
		t.Fatalf("expected exit status to be not ok for nonexistent device json")
	}
}

func readTestJSON(t *testing.T, path string) gjson.Result {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read json: %v", err)
	}
	return gjson.ParseBytes(data)
}
