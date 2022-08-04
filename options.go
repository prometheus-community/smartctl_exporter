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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var (
	exporterVersion = "0.5"
)

// SMARTOptions is a inner representation of a options
type SMARTOptions struct {
	BindTo                string `yaml:"bind_to"`
	URLPath               string `yaml:"url_path"`
	FakeJSON              bool   `yaml:"fake_json"`
	SMARTctlLocation      string `yaml:"smartctl_location"`
	CollectPeriod         string `yaml:"collect_not_more_than_period"`
	CollectPeriodDuration time.Duration
	Devices               []string `yaml:"devices"`
}

// Options is a representation of a options
type Options struct {
	SMARTctl SMARTOptions `yaml:"smartctl_exporter"`
}

// Parse options from yaml config file
func loadOptions() Options {
	configFile := flag.String("config", "/etc/smartctl_exporter.yaml", "Path to smartctl_exporter config file")
	verbose := flag.Bool("verbose", false, "Verbose log output")
	debug := flag.Bool("debug", false, "Debug log output")
	version := flag.Bool("version", false, "Show application version and exit")
	flag.Parse()

	if *version {
		fmt.Printf("smartctl_exporter version: %s\n", exporterVersion)
		os.Exit(0)
	}

	logger = newLogger(*verbose, *debug)

	logger.Verbose("Read options from %s\n", *configFile)
	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		logger.Panic("Failed read %s: %s", configFile, err)
	}

	opts := Options{
		SMARTOptions{
			BindTo:           "9633",
			URLPath:          "/metrics",
			FakeJSON:         false,
			SMARTctlLocation: "/usr/sbin/smartctl",
			CollectPeriod:    "60s",
			Devices:          []string{},
		},
	}

	if yaml.Unmarshal(yamlFile, &opts) != nil {
		logger.Panic("Failed parse %s: %s", configFile, err)
	}

	d, err := time.ParseDuration(opts.SMARTctl.CollectPeriod)
	if err != nil {
		logger.Panic("Failed read collect_not_more_than_period (%s): %s", opts.SMARTctl.CollectPeriod, err)
	}

	opts.SMARTctl.CollectPeriodDuration = d

	logger.Debug("Parsed options: %s", opts)
	return opts
}
