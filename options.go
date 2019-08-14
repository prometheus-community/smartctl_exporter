package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var (
	exporterVersion = "0.5"
)

// SMARTOptions is a inner representation of a options
type SMARTOptions struct {
	BindTo           string   `yaml:"bind_to"`
	URLPath          string   `yaml:"url_path"`
	FakeJSON         bool     `yaml:"fake_json"`
	SMARTctlLocation string   `yaml:"smartctl_location"`
	Devices          []string `yaml:"devices"`
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
			BindTo:           "9631",
			URLPath:          "/metrics",
			FakeJSON:         false,
			SMARTctlLocation: "/usr/sbin/smartctl",
			Devices:          []string{},
		},
	}

	if yaml.Unmarshal(yamlFile, &opts) != nil {
		logger.Panic("Failed parse %s: %s", configFile, err)
	}
	logger.Debug("Parsed options: %s", opts)
	return opts
}
