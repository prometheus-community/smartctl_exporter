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
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	SmartctlPath            *string  `yaml:"smartctl.path"`
	SmartctlInterval        *string  `yaml:"smartctl.interval"`
	SmartctlRescan          *string  `yaml:"smartctl.rescan"`
	SmartctlScan            *bool    `yaml:"smartctl.scan"`
	SmartctlDevices         []string `yaml:"smartctl.device"`
	SmartctlDeviceExclude   *string  `yaml:"smartctl.device-exclude"`
	SmartctlDeviceInclude   *string  `yaml:"smartctl.device-include"`
	SmartctlScanDeviceTypes []string `yaml:"smartctl.scan-device-type"`
	SmartctlPowerModeCheck  *string  `yaml:"smartctl.powermode-check"`
	WebListenAddress        *string  `yaml:"web.listen-address"`
	WebTelemetryPath        *string  `yaml:"web.telemetry-path"`
}

func findConfigFile(args []string) string {
	for i, arg := range args {
		if arg == "--config.file" && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, "--config.file=") {
			return strings.TrimPrefix(arg, "--config.file=")
		}
	}
	return ""
}

func loadConfigFile(path string) (Config, error) {
	if path == "" {
		return Config{}, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file %q: %w", path, err)
	}
	defer file.Close()

	var cfg Config
	scanner := bufio.NewScanner(file)
	var listKey string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "- ") {
			if listKey == "" {
				return Config{}, fmt.Errorf("unexpected list item in %q: %s", path, line)
			}
			item := trimConfigValue(strings.TrimSpace(strings.TrimPrefix(line, "-")))
			if item == "" {
				continue
			}
			if err := appendListValue(&cfg, listKey, item); err != nil {
				return Config{}, fmt.Errorf("parse config file %q: %w", path, err)
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("invalid config line in %q: %s", path, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if value == "" {
			listKey = key
			continue
		}
		listKey = ""
		value = trimConfigValue(value)

		if err := assignConfigValue(&cfg, key, value); err != nil {
			return Config{}, fmt.Errorf("parse config file %q: %w", path, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return Config{}, fmt.Errorf("read config file %q: %w", path, err)
	}
	return cfg, nil
}

func stringValueOrDefault(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}

func trimConfigValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "\"")
	value = strings.Trim(value, "'")
	value = strings.ReplaceAll(value, `\\`, `\`)
	return value
}

func assignConfigValue(cfg *Config, key, value string) error {
	switch key {
	case "smartctl.path":
		cfg.SmartctlPath = &value
	case "smartctl.interval":
		cfg.SmartctlInterval = &value
	case "smartctl.rescan":
		cfg.SmartctlRescan = &value
	case "smartctl.scan":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid smartctl.scan value %q", value)
		}
		cfg.SmartctlScan = &parsed
	case "smartctl.device":
		cfg.SmartctlDevices = append(cfg.SmartctlDevices, value)
	case "smartctl.device-exclude":
		cfg.SmartctlDeviceExclude = &value
	case "smartctl.device-include":
		cfg.SmartctlDeviceInclude = &value
	case "smartctl.scan-device-type":
		cfg.SmartctlScanDeviceTypes = append(cfg.SmartctlScanDeviceTypes, value)
	case "smartctl.powermode-check":
		cfg.SmartctlPowerModeCheck = &value
	case "web.listen-address":
		cfg.WebListenAddress = &value
	case "web.telemetry-path":
		cfg.WebTelemetryPath = &value
	default:
		return fmt.Errorf("unsupported config key %q", key)
	}
	return nil
}

func appendListValue(cfg *Config, key, value string) error {
	switch key {
	case "smartctl.device":
		cfg.SmartctlDevices = append(cfg.SmartctlDevices, value)
	case "smartctl.scan-device-type":
		cfg.SmartctlScanDeviceTypes = append(cfg.SmartctlScanDeviceTypes, value)
	default:
		return fmt.Errorf("unsupported list key %q", key)
	}
	return nil
}
