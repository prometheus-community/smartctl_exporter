package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var (
	exporterVersion = "v0.1"
	OSCommand string
)

// SMARTOptions is a inner representation of a options
type SMARTOptions struct {
	BindTo                string `yaml:"bind_to"`
	URLPath               string `yaml:"url_path"`
	FakeJSON              bool   `yaml:"fake_json"`
	//SMARTctlLocation      string `yaml:"smartctl_location"`
	CollectPeriod         string `yaml:"collect_not_more_than_period"`
	CollectPeriodDuration time.Duration
	Devices               []string `yaml:"devices"`
}

// Options is a representation of a options
type Options struct {
	SMARTctl SMARTOptions `yaml:"smartctl_exporter"`
}

func GetOSCommand()  {
	cmd,err:=exec.Command("which","smartctl").Output()
	if err!=nil {
		log.Panic("系统为安装命令")
	}


	OSCommand=strings.Split(string(cmd),"\n")[0]
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
			//SMARTctlLocation: OSCommand,
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
