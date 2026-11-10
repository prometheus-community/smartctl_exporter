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
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestResultCodeIsOkQuietExitStatusBits(t *testing.T) {
	saved := *smartctlQuietExitStatusBits
	t.Cleanup(func() { *smartctlQuietExitStatusBits = saved })

	tests := []struct {
		name       string
		quietBits  []int
		exitStatus int64
		wantOk     bool
		wantWarn   []string
		wantDebug  []string
	}{
		{
			name:       "error log records are a warning by default",
			exitStatus: 1 << 6,
			wantOk:     true,
			wantWarn:   []string{"The device error log contains records of errors"},
		},
		{
			name:       "quiet bit is logged at debug level",
			quietBits:  []int{6},
			exitStatus: 1<<6 | 1<<4,
			wantOk:     true,
			wantWarn:   []string{"We found prefail Attributes <= threshold"},
			wantDebug:  []string{"The device error log contains records of errors"},
		},
		{
			name:       "failed device open is still an error",
			quietBits:  []int{6},
			exitStatus: 1<<1 | 1<<6,
			wantOk:     false,
			wantDebug:  []string{"The device error log contains records of errors"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			*smartctlQuietExitStatusBits = test.quietBits
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

			if ok := resultCodeIsOk(logger, Device{Name: "/dev/sda", Type: "auto"}, test.exitStatus); ok != test.wantOk {
				t.Errorf("resultCodeIsOk() = %v, want %v", ok, test.wantOk)
			}

			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			for _, msg := range test.wantWarn {
				if !hasLogLine(lines, "level=WARN", msg) {
					t.Errorf("expected a warning %q, got:\n%s", msg, buf.String())
				}
			}
			for _, msg := range test.wantDebug {
				if !hasLogLine(lines, "level=DEBUG", msg) {
					t.Errorf("expected a debug message %q, got:\n%s", msg, buf.String())
				}
			}
		})
	}
}

func hasLogLine(lines []string, level, msg string) bool {
	for _, line := range lines {
		if strings.Contains(line, level) && strings.Contains(line, msg) {
			return true
		}
	}
	return false
}

func TestValidateQuietExitStatusBits(t *testing.T) {
	if err := validateQuietExitStatusBits([]int{2, 6, 7}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	for _, bit := range []int{0, 1, 8} {
		if err := validateQuietExitStatusBits([]int{bit}); err == nil {
			t.Errorf("expected an error for bit %d", bit)
		}
	}
}
