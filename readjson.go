package main

import (
	"fmt"
	"io/ioutil"
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
func readSMARTctl(device string) gjson.Result {
	logger.Debug("Collecting S.M.A.R.T. counters, device: %s", device)
	out, err := exec.Command(OSCommand, "--json", "--xall", device).Output()
	if err != nil {
		logger.Warning("S.M.A.R.T. output reading error: %s", err)
	}
	return parseJSON(string(out))
}

// Select json source and parse
func readData(device string) gjson.Result {
	if options.SMARTctl.FakeJSON {
		return readFakeSMARTctl(device)
	}

	if value, ok := jsonCache[device]; ok {
		// logger.Debug("Cache exists")
		if time.Now().After(value.LastCollect.Add(options.SMARTctl.CollectPeriodDuration)) {
			// logger.Debug("Cache update")
			jsonCache[device] = JSONCache{JSON: readSMARTctl(device), LastCollect: time.Now()}
		}
	} else {
		// logger.Debug("Cache not exists")
		jsonCache[device] = JSONCache{JSON: readSMARTctl(device), LastCollect: time.Now()}
	}
	return jsonCache[device].JSON
}
