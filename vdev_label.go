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
	"os/exec"
	"strings"
)

func newDevice(logger *slog.Logger, name, deviceType string) Device {
	return newDeviceWithVdevLookup(logger, name, deviceType, *smartctlVdevLabel, lookupVdevLabel)
}

func newDeviceWithVdevLookup(logger *slog.Logger, name, deviceType string, enabled bool, lookup func(string) (string, error)) Device {
	device := Device{
		Name:  name,
		Type:  deviceType,
		Label: buildDeviceLabel(name, deviceType),
	}
	if !enabled {
		return device
	}

	label, err := lookup(name)
	if err != nil {
		logger.Debug("Unable to query udev properties", "device", name, "err", err)
		return device
	}
	if label != "" {
		device.Label = label
		logger.Debug("Using ID_VDEV as device label", "device", name, "label", label)
	}
	return device
}

func lookupVdevLabel(name string) (string, error) {
	output, err := exec.Command("udevadm", "info", "--query=property", name).Output()
	if err != nil {
		return "", err
	}
	return parseVdevLabel(string(output)), nil
}

func parseVdevLabel(output string) string {
	for line := range strings.Lines(output) {
		key, value, found := strings.Cut(strings.TrimSpace(line), "=")
		if found && key == "ID_VDEV" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
