package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/tidwall/gjson"
)

// Parse json to gjson object
func parseJSON(data string) (gjson.Result, error) {
	if !gjson.Valid(data) {
		return gjson.Parse("{}"), errors.New("Invalid JSON")
	}
	return gjson.Parse(data), nil
}

// Reading fake smartctl json
func readFakeSMARTctl(device string) (gjson.Result, error) {
	splitted := strings.Split(device, "/")
	filename := fmt.Sprintf("%s.json", splitted[len(splitted)-1])
	logger.Verbose("Read fake S.M.A.R.T. data from json: %s", filename)
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Fake S.M.A.R.T. data reading error: %s", err)
		return parseJSON("{}")
	}
	return parseJSON(string(jsonFile))
}

// Get json from smartctl and parse it
func readSMARTctl(device string) (gjson.Result, error) {
	logger.Debug("Collecting S.M.A.R.T. counters, device: %s...", device)
	out, err := exec.Command(options.SMARTctl.SMARTctlLocation, "--json", "--xall", device).Output()
	if err != nil {
		logger.Error("S.M.A.R.T. output reading error: %s", err)
	}
	return parseJSON(string(out))
}

// Select json source and parse
func readData(device string) (gjson.Result, error) {
	if options.SMARTctl.FakeJSON {
		return readFakeSMARTctl(device)
	}
	return readSMARTctl(device)
}
