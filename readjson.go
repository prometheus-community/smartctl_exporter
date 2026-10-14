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
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

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
func readFakeSMARTctl(logger *slog.Logger, device Device) gjson.Result {
	s := strings.Split(device.Name, "/")
	filename := fmt.Sprintf("debug/%s.json", s[len(s)-1])
	logger.Debug("Read fake S.M.A.R.T. data from json", "filename", filename)
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("Fake S.M.A.R.T. data reading error", "err", err)
		return parseJSON("{}")
	}
	return parseJSON(string(jsonFile))
}

var errSmartctlTimeout = errors.New("smartctl timed out")

// runSmartctl runs smartctl and returns its stdout. If --smartctl.timeout is
// set and elapses, the process is killed and errSmartctlTimeout is returned
// without waiting for it to exit: a process stuck in uninterruptible I/O (D
// state) ignores SIGKILL, and waiting for it would block forever.
// onExit, if not nil, is called once the process has really exited (or failed
// to start).
func runSmartctl(args []string, onExit func()) ([]byte, error) {
	cmd := exec.Command(*smartctlPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Start(); err != nil {
		if onExit != nil {
			onExit()
		}
		return nil, err
	}
	done := make(chan error, 1)
	go func() {
		err := cmd.Wait()
		if onExit != nil {
			onExit()
		}
		done <- err
	}()

	var timeout <-chan time.Time // nil (blocks forever) when no timeout is set
	if *smartctlTimeout > 0 {
		timer := time.NewTimer(*smartctlTimeout)
		defer timer.Stop()
		timeout = timer.C
	}
	select {
	case err := <-done:
		return stdout.Bytes(), err
	case <-timeout:
		_ = cmd.Process.Kill()
		return nil, errSmartctlTimeout
	}
}

// Devices with a smartctl process that has not exited yet.
var inFlight sync.Map

// Result of the last read of each device: true if ok, false if it failed.
var deviceScrapeOk sync.Map

// scrapeFailed drops the cached data of a device so that no stale metrics are
// exported for it. The cache is populated again by the next successful read.
func scrapeFailed(device Device) {
	jsonCache.Delete(device)
	deviceScrapeOk.Store(device, false)
}

// Get json from smartctl and parse it
func readSMARTctl(logger *slog.Logger, device Device, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	var smartctlArgs = []string{"--json", "--info", "--health", "--attributes", "--capabilities", "--tolerance=verypermissive", "--nocheck=" + *smartctlPowerModeCheck, "--format=brief", "--log=error", "--device=" + device.Type, device.Name}

	logger.Debug("Calling smartctl with args", "args", strings.Join(smartctlArgs, " "))
	out, err := runSmartctl(smartctlArgs, func() { inFlight.Delete(device) })
	if err != nil {
		logger.Warn("S.M.A.R.T. output reading", "err", err, "device", device)
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			// Timeout or smartctl could not run: there is no output to use.
			scrapeFailed(device)
			return
		}
	}
	// Accommodate a smartmontools pre-7.3 bug
	json := parseJSON(strings.TrimPrefix(string(out), "  Pending defect count:"))
	rcOk := resultCodeIsOk(logger, device, json.Get("smartctl.exit_status").Int())
	jsonOk := jsonIsOk(logger, json)
	logger.Debug("Collected S.M.A.R.T. json data", "device", device, "duration", time.Since(start))
	if rcOk && jsonOk {
		jsonCache.Store(device, JSONCache{JSON: json, LastCollect: time.Now()})
		deviceScrapeOk.Store(device, true)
	} else if !jsonOk {
		scrapeFailed(device)
	}
	// rc not ok without a json error is a device in standby/low-power mode: keep
	// the previous state as before.
}

func readSMARTctlDevices(logger *slog.Logger) gjson.Result {
	logger.Debug("Scanning for devices")
	scanArgs := []string{"--json", "--scan"}
	for _, d := range *smartctlScanDeviceTypes {
		scanArgs = append(scanArgs, "--device", d)
	}
	out, err := runSmartctl(scanArgs, nil)
	if exiterr, ok := err.(*exec.ExitError); ok {
		logger.Debug("Exit Status", "exit_code", exiterr.ExitCode())
		// The smartctl command returns 2 if devices are sleeping, ignore this error.
		if exiterr.ExitCode() != 2 {
			logger.Warn("S.M.A.R.T. output reading error", "err", err)
			return gjson.Result{}
		}
	} else if err != nil {
		logger.Warn("S.M.A.R.T. output reading error", "err", err)
		return gjson.Result{}
	}
	return parseJSON(string(out))
}

// Refresh all devices' json
func refreshAllDevices(logger *slog.Logger, devices []Device) {
	if *smartctlFakeData {
		return
	}

	var wg sync.WaitGroup
	for _, device := range devices {
		cacheValue, cacheOk := jsonCache.Load(device)
		if !cacheOk || time.Now().After(cacheValue.(JSONCache).LastCollect.Add(*smartctlInterval)) {
			// Do not start another smartctl while the previous one is still stuck.
			if _, running := inFlight.LoadOrStore(device, struct{}{}); running {
				logger.Debug("smartctl still running, skipping", "device", device)
				continue
			}
			wg.Add(1)
			go readSMARTctl(logger, device, &wg)
		}
	}
	wg.Wait()
}

func readData(logger *slog.Logger, device Device) gjson.Result {
	if *smartctlFakeData {
		return readFakeSMARTctl(logger, device)
	}

	cacheValue, found := jsonCache.Load(device)
	if !found {
		// A failed read was already logged and its cache dropped.
		if ok, failed := deviceScrapeOk.Load(device); failed && !ok.(bool) {
			logger.Debug("no data for device, last read failed", "device", device)
		} else {
			logger.Warn("device not found", "device", device)
		}
		return gjson.Result{}
	}
	return cacheValue.(JSONCache).JSON
}

// Parse smartctl return code
func resultCodeIsOk(logger *slog.Logger, device Device, SMARTCtlResult int64) bool {
	result := true
	if SMARTCtlResult > 0 {
		b := SMARTCtlResult
		if (b & 1) != 0 {
			logger.Error("Command line did not parse", "device", device)
			result = false
		}
		if (b & (1 << 1)) != 0 {
			logger.Error("Device open failed, device did not return an IDENTIFY DEVICE structure, or device is in a low-power mode", "device", device)
			result = false
		}
		if (b & (1 << 2)) != 0 {
			logger.Warn("Some SMART or other ATA command to the disk failed, or there was a checksum error in a SMART data structure", "device", device)
		}
		if (b & (1 << 3)) != 0 {
			logger.Warn("SMART status check returned 'DISK FAILING'", "device", device)
		}
		if (b & (1 << 4)) != 0 {
			logger.Warn("We found prefail Attributes <= threshold", "device", device)
		}
		if (b & (1 << 5)) != 0 {
			logger.Warn("SMART status check returned 'DISK OK' but we found that some (usage or prefail) Attributes have been <= threshold at some time in the past", "device", device)
		}
		if (b & (1 << 6)) != 0 {
			logger.Warn("The device error log contains records of errors", "device", device)
		}
		if (b & (1 << 7)) != 0 {
			logger.Warn("The device self-test log contains records of errors. [ATA only] Failed self-tests outdated by a newer successful extended self-test are ignored", "device", device)
		}
	}
	return result
}

// Check json
func jsonIsOk(logger *slog.Logger, json gjson.Result) bool {
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
