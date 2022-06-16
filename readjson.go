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
	s := strings.Split(device, "/")
	filename := fmt.Sprintf("debug/%s.json", s[len(s)-1])
	logger.Verbose("Read fake S.M.A.R.T. data from json: %s", filename)
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Fake S.M.A.R.T. data reading error: %s", err)
		return parseJSON("{}")
	}
	return parseJSON(string(jsonFile))
}

// Get json from smartctl and parse it
func readSMARTctl(device string) (gjson.Result, bool) {
	logger.Debug("Collecting S.M.A.R.T. counters, device: %s", device)
	out, err := exec.Command(options.SMARTctl.SMARTctlLocation, "--json", "--xall", device).Output()
	if err != nil {
		logger.Warning("S.M.A.R.T. output reading error: %s", err)
	}
	json := parseJSON(string(out))
	rcOk := resultCodeIsOk(json.Get("smartctl.exit_status").Int())
	jsonOk := jsonIsOk(json)
	return json, rcOk && jsonOk
}

func readSMARTctlDevices() gjson.Result {
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
		timeToScan := false
		if cacheOk {
			timeToScan = time.Now().After(cacheValue.LastCollect.Add(options.SMARTctl.CollectPeriodDuration))
		} else {
			timeToScan = true
		}

		if timeToScan {
			json, ok := readSMARTctl(device)
			if ok {
				jsonCache[device] = JSONCache{JSON: json, LastCollect: time.Now()}
				return jsonCache[device].JSON, nil
			}
			return gjson.Parse("{}"), fmt.Errorf("smartctl returned bad data for device %s", device)
		}
		return gjson.Parse("{}"), fmt.Errorf("Too early collect called for device %s", device)
	}
	return gjson.Parse("{}"), fmt.Errorf("Device %s unavialable", device)
}

// Parse smartctl return code
func resultCodeIsOk(SMARTCtlResult int64) bool {
	result := true
	if SMARTCtlResult > 0 {
		bits := fmt.Sprintf("%08b", SMARTCtlResult)
		// logger.Debug("Return code: %d: %s", SMARTCtlResult, bits)
		if bits[0] == '1' {
			logger.Error("Command line did not parse.")
			result = false
		}
		if bits[1] == '1' {
			logger.Error("Device open failed, device did not return an IDENTIFY DEVICE structure, or device is in a low-power mode")
			result = false
		}
		if bits[2] == '1' {
			logger.Warning("Some SMART or other ATA command to the disk failed, or there was a checksum error in a SMART data structure")
		}
		if bits[3] == '1' {
			logger.Warning("SMART status check returned 'DISK FAILING'.")
		}
		if bits[4] == '1' {
			logger.Warning("We found prefail Attributes <= threshold.")
		}
		if bits[5] == '1' {
			logger.Warning("SMART status check returned 'DISK OK' but we found that some (usage or prefail) Attributes have been <= threshold at some time in the past.")
		}
		if bits[6] == '1' {
			logger.Warning("The device error log contains records of errors.")
		}
		if bits[7] == '1' {
			logger.Warning("The device self-test log contains records of errors. [ATA only] Failed self-tests outdated by a newer successful extended self-test are ignored.")
		}
	}
	return result
}

// Check json
func jsonIsOk(json gjson.Result) bool {
	messages := json.Get("smartctl.messages")
	// logger.Debug(messages.String())
	if messages.Exists() {
		for _, message := range messages.Array() {
			if message.Get("severity").String() == "error" {
				logger.Error(message.Get("string").String())
				return false
			}
		}
	}
	return true
}
