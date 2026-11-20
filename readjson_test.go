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
	"testing"
)

func TestJSONIsOk(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	errorLogFailed := `{"string": "Read 16 entries from Error Information Log failed: PRP Offset Invalid (0x013)", "severity": "error"}`

	tests := []struct {
		name string
		json string
		want bool
	}{
		{
			name: "no messages",
			json: `{"smart_status": {"passed": true}}`,
			want: true,
		},
		{
			name: "warning only",
			json: `{"smartctl": {"messages": [{"string": "Warning: ATA error count 3", "severity": "warning"}]}}`,
			want: true,
		},
		{
			name: "error with S.M.A.R.T. status",
			json: `{"smartctl": {"messages": [` + errorLogFailed + `], "exit_status": 4}, "smart_status": {"passed": true}}`,
			want: true,
		},
		{
			name: "error without S.M.A.R.T. status",
			json: `{"smartctl": {"messages": [{"string": "Smartctl open device: /dev/sdz failed: No such device", "severity": "error"}], "exit_status": 2}}`,
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := jsonIsOk(logger, parseJSON(test.json)); got != test.want {
				t.Errorf("jsonIsOk() = %v, want %v", got, test.want)
			}
		})
	}
}
