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
	"errors"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

func TestParseVdevLabel(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "ID_VDEV property",
			output:   "DEVNAME=/dev/sda\nID_MODEL=ST12000NM000J\nID_VDEV=2-1\n",
			expected: "2-1",
		},
		{
			name:     "value containing equals",
			output:   "ID_VDEV=pool=slot-1\n",
			expected: "pool=slot-1",
		},
		{
			name:     "empty property",
			output:   "ID_VDEV=\n",
			expected: "",
		},
		{
			name:     "property absent",
			output:   "DEVNAME=/dev/sda\nID_MODEL=ST12000NM000J\n",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := parseVdevLabel(test.output); actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestNewDeviceWithoutVdevLabel(t *testing.T) {
	device := newDeviceWithVdevLookup(slog.Default(), "/dev/bus/0", "megaraid,1", false, func(string) (string, error) {
		t.Fatal("lookup called when VDEV labels are disabled")
		return "", nil
	})
	if device.Label != "bus_0_megaraid_1" {
		t.Errorf("expected default device label, got %q", device.Label)
	}
}

func TestNewDeviceWithVdevLabel(t *testing.T) {
	device := newDeviceWithVdevLookup(slog.Default(), "/dev/sda", "auto", true, func(name string) (string, error) {
		if name != "/dev/sda" {
			t.Errorf("expected lookup for /dev/sda, got %q", name)
		}
		return "2-1", nil
	})
	if device.Label != "2-1" {
		t.Errorf("expected VDEV label, got %q", device.Label)
	}
}

func TestNewDeviceVdevLabelFallback(t *testing.T) {
	tests := []struct {
		name   string
		lookup func(string) (string, error)
	}{
		{
			name: "property absent",
			lookup: func(string) (string, error) {
				return "", nil
			},
		},
		{
			name: "lookup error",
			lookup: func(string) (string, error) {
				return "", errors.New("udevadm unavailable")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			device := newDeviceWithVdevLookup(slog.Default(), "/dev/sda", "auto", true, test.lookup)
			if device.Label != "sda" {
				t.Errorf("expected default device label, got %q", device.Label)
			}
		})
	}
}

func TestNewSMARTctlUsesConfiguredDeviceLabel(t *testing.T) {
	json := gjson.Parse(`{
		"device": {"name": "/dev/sda", "type": "sat", "protocol": "ATA"},
		"model_name": "ST12000NM000J"
	}`)
	ch := make(chan prometheus.Metric)
	device := Device{Name: "/dev/sda", Type: "sat", Label: "2-1"}

	smart := NewSMARTctl(slog.Default(), json, ch, device)
	if smart.device.device != "2-1" {
		t.Errorf("expected configured device label, got %q", smart.device.device)
	}
}
