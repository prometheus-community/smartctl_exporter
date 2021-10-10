package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// JSONCache caching json
type JSONCache struct {
	JSON        gjson.Result
	LastCollect time.Time
}

var (
	jsonCache map[string]JSONCache
)

func init() {
	jsonCache = make(map[string]JSONCache)
}

// Parse json to gjson object
func parseJSON(data string) gjson.Result {
	if !gjson.Valid(data) {
		return gjson.Parse("{}")
	}
	return gjson.Parse(data)
}

// Reading fake smartctl json
func readFakeSMARTctl(device string) gjson.Result {
	splitted := strings.Split(device, "/")
	filename := fmt.Sprintf("debug/%s.json", splitted[len(splitted)-1])
	logger.Verbose("Read fake S.M.A.R.T. data from json: %s", filename)
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Fake S.M.A.R.T. data reading error: %s", err)
		return parseJSON("{}")
	}
	return parseJSON(string(jsonFile))
}

// Get json from smartctl and parse it
func readSMARTctlDevices(device string) gjson.Result {
	logger.Debug("Collecting S.M.A.R.T. counters, device: %s", device)
	out, err := exec.Command(options.SMARTctl.SMARTctlLocation, "--json", "--xall", device).Output()
	if err != nil {
		logger.Warning("S.M.A.R.T. output reading error: %s", err)
	}
	return parseJSON(string(out))
}

func scanSMARTctlDevices() gjson.Result {
	logger.Debug("Collecting devices")
	out, err := exec.Command(options.SMARTctl.SMARTctlLocation, "--json", "--scan-open").Output()
	if err != nil {
		logger.Warning("S.M.A.R.T. output reading error: %s", err)
	}
	return parseJSON(string(out))
}

// Select json source and parse
func readData(device string) (gjson.Result, error) {
	if options.SMARTctl.FakeJSON {
		return readFakeSMARTctl(device), nil
	}

	if _, err := os.Stat(device); err == nil {
		cacheValue, cacheOk := jsonCache[device]
		if !cacheOk || time.Now().After(cacheValue.LastCollect.Add(options.SMARTctl.CollectPeriodDuration)) {
			jsonCache[device] = JSONCache{JSON: readSMARTctlDevices(device), LastCollect: time.Now()}
			return jsonCache[device].JSON, nil
		}
		return gjson.Parse("{}"), fmt.Errorf("Too early collect called for device %s", device)
	}
	return gjson.Parse("{}"), fmt.Errorf("Device %s unavialable", device)
}

// Check json
func jsonIsOk(json gjson.Result) {
	messages := json.Get("smartctl.messages")
	if messages.Exists() {
		for _, message := range messages.Array() {
			if message.Get("severity").String() == "error" {
				logger.Error(message.Get("string").String())
			}
		}
	}
}
