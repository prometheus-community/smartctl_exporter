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
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/tidwall/gjson"
)

// JSONCache caching json
type JSONCache struct {
	JSON        gjson.Result
	LastCollect time.Time
}

var (
	jsonCache sync.Map
)

func init() {
	jsonCache.Store("", JSONCache{})
}

// Parse json to gjson object
func parseJSON(data string) gjson.Result {
	if !gjson.Valid(data) {
		return gjson.Parse("{}")
	}
	return gjson.Parse(data)
}

// Reading fake smartctl json
func readFakeSMARTctl(logger log.Logger, device Device) gjson.Result {
	s := strings.Split(device.Name, "/")
	filename := fmt.Sprintf("debug/%s.json", s[len(s)-1])
	level.Debug(logger).Log("msg", "Read fake S.M.A.R.T. data from json", "filename", filename)
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		level.Error(logger).Log("msg", "Fake S.M.A.R.T. data reading error", "err", err)
		return parseJSON("{}")
	}
	return parseJSON(string(jsonFile))
}

// Get json from smartctl and parse it
func readSMARTctl(logger log.Logger, device Device) (gjson.Result, bool) {
	start := time.Now()
	out, err := exec.Command(*smartctlPath, "--json", "--info", "--health", "--attributes", "--tolerance=verypermissive", "--nocheck=standby", "--format=brief", "--log=error", "--device="+device.Type, device.Name).Output()
	if err != nil {
		level.Warn(logger).Log("msg", "S.M.A.R.T. output reading", "err", err, "device", device.Info_Name)
	}
	json := parseJSON(string(out))
	rcOk := resultCodeIsOk(logger, device, json.Get("smartctl.exit_status").Int())
	jsonOk := jsonIsOk(logger, json)
	level.Debug(logger).Log("msg", "Collected S.M.A.R.T. json data", "device", device.Info_Name, "duration", time.Since(start))
	return json, rcOk && jsonOk
}

func readSMARTctlDevices(logger log.Logger) gjson.Result {
	level.Debug(logger).Log("msg", "Scanning for devices")
	out, err := exec.Command(*smartctlPath, "--json", "--scan").Output()
	if exiterr, ok := err.(*exec.ExitError); ok {
		level.Debug(logger).Log("msg", "Exit Status", "exit_code", exiterr.ExitCode())
		// The smartctl command returns 2 if devices are sleeping, ignore this error.
		if exiterr.ExitCode() != 2 {
			level.Warn(logger).Log("msg", "S.M.A.R.T. output reading error", "err", err)
			return gjson.Result{}
		}
	}
	return parseJSON(string(out))
}

// Select json source and parse
func readData(logger log.Logger, device Device) gjson.Result {
	if *smartctlFakeData {
		return readFakeSMARTctl(logger, device)
	}

	cacheValue, cacheOk := jsonCache.Load(device)
	if !cacheOk || time.Now().After(cacheValue.(JSONCache).LastCollect.Add(*smartctlInterval)) {
		json, ok := readSMARTctl(logger, device)
		if ok {
			jsonCache.Store(device, JSONCache{JSON: json, LastCollect: time.Now()})
			j, found := jsonCache.Load(device)
			if !found {
				level.Warn(logger).Log("msg", "device not found", "device", device.Info_Name)
			}
			return j.(JSONCache).JSON
		}
		return gjson.Result{}
	}
	return cacheValue.(JSONCache).JSON
}

// Parse smartctl return code
func resultCodeIsOk(logger log.Logger, device Device, SMARTCtlResult int64) bool {
	result := true
	if SMARTCtlResult > 0 {
		b := SMARTCtlResult
		if (b & 1) != 0 {
			level.Error(logger).Log("msg", "Command line did not parse", "device", device.Info_Name)
			result = false
		}
		if (b & (1 << 1)) != 0 {
			level.Error(logger).Log("msg", "Device open failed, device did not return an IDENTIFY DEVICE structure, or device is in a low-power mode", "device", device.Info_Name)
			result = false
		}
		if (b & (1 << 2)) != 0 {
			level.Warn(logger).Log("msg", "Some SMART or other ATA command to the disk failed, or there was a checksum error in a SMART data structure", "device", device.Info_Name)
		}
		if (b & (1 << 3)) != 0 {
			level.Warn(logger).Log("msg", "SMART status check returned 'DISK FAILING'", "device", device.Info_Name)
		}
		if (b & (1 << 4)) != 0 {
			level.Warn(logger).Log("msg", "We found prefail Attributes <= threshold", "device", device.Info_Name)
		}
		if (b & (1 << 5)) != 0 {
			level.Warn(logger).Log("msg", "SMART status check returned 'DISK OK' but we found that some (usage or prefail) Attributes have been <= threshold at some time in the past", "device", device.Info_Name)
		}
		if (b & (1 << 6)) != 0 {
			level.Warn(logger).Log("msg", "The device error log contains records of errors", "device", device.Info_Name)
		}
		if (b & (1 << 7)) != 0 {
			level.Warn(logger).Log("msg", "The device self-test log contains records of errors. [ATA only] Failed self-tests outdated by a newer successful extended self-test are ignored", "device", device.Info_Name)
		}
	}
	return result
}

// Check json
func jsonIsOk(logger log.Logger, json gjson.Result) bool {
	messages := json.Get("smartctl.messages")
	// logger.Debug(messages.String())
	if messages.Exists() {
		for _, message := range messages.Array() {
			if message.Get("severity").String() == "error" {
				level.Error(logger).Log("msg", message.Get("string").String())
				return false
			}
		}
	}
	return true
}
