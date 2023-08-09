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
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

// SMARTDevice - short info about device
type SMARTDevice struct {
	device string
	serial string
	family string
	model  string
}

// SMARTctl object
type SMARTctl struct {
	ch     chan<- prometheus.Metric
	json   gjson.Result
	logger log.Logger
	device SMARTDevice
}

// NewSMARTctl is smartctl constructor
func NewSMARTctl(logger log.Logger, json gjson.Result, ch chan<- prometheus.Metric) SMARTctl {
	return SMARTctl{
		ch:     ch,
		json:   json,
		logger: logger,
		device: SMARTDevice{
			device: strings.TrimPrefix(strings.TrimSpace(json.Get("device.name").String()), "/dev/"),
			serial: strings.TrimSpace(json.Get("serial_number").String()),
			family: strings.TrimSpace(json.Get("model_family").String()),
			model:  strings.TrimSpace(json.Get("model_name").String()),
		},
	}
}

// Collect metrics
func (smart *SMARTctl) Collect() {
	level.Debug(smart.logger).Log("msg", "Collecting metrics from", "device", smart.device.device, "family", smart.device.family, "model", smart.device.model)
	smart.mineExitStatus()
	smart.mineDevice()
	smart.mineCapacity()
	smart.mineInterfaceSpeed()
	smart.mineDeviceAttribute()
	smart.minePowerOnSeconds()
	smart.mineRotationRate()
	smart.mineTemperatures()
	smart.minePowerCycleCount()
	smart.mineDeviceSCTStatus()
	smart.mineDeviceStatistics()
	smart.mineDeviceStatus()
	smart.mineDeviceErrorLog()
	smart.mineDeviceSelfTestLog()
	smart.mineDeviceERC()
	smart.minePercentageUsed()
	smart.mineAvailableSpare()
	smart.mineAvailableSpareThreshold()
	smart.mineCriticalWarning()
	smart.mineMediaErrors()
	smart.mineNumErrLogEntries()
	smart.mineBytesRead()
	smart.mineBytesWritten()
	smart.mineSmartStatus()
	smart.mineSCSIGrownDefectList()
	smart.mineSCSIErrorCounterLog()
}

func (smart *SMARTctl) mineExitStatus() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceExitStatus,
		prometheus.GaugeValue,
		smart.json.Get("smartctl.exit_status").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineDevice() {
	device := smart.json.Get("device")
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceModel,
		prometheus.GaugeValue,
		1,
		smart.device.device,
		device.Get("type").String(),
		device.Get("protocol").String(),
		smart.device.family,
		smart.device.model,
		smart.device.serial,
		GetStringIfExists(smart.json, "ata_additional_product_id", "unknown"),
		smart.json.Get("firmware_version").String(),
		smart.json.Get("ata_version.string").String(),
		smart.json.Get("sata_version.string").String(),
		smart.json.Get("form_factor.name").String(),
	)
}

func (smart *SMARTctl) mineCapacity() {
	capacity := smart.json.Get("user_capacity")
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceCapacityBlocks,
		prometheus.GaugeValue,
		capacity.Get("blocks").Float(),
		smart.device.device,
	)
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceCapacityBytes,
		prometheus.GaugeValue,
		capacity.Get("bytes").Float(),
		smart.device.device,
	)
	for _, blockType := range []string{"logical", "physical"} {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceBlockSize,
			prometheus.GaugeValue,
			smart.json.Get(fmt.Sprintf("%s_block_size", blockType)).Float(),
			smart.device.device,
			blockType,
		)
	}
}

func (smart *SMARTctl) mineInterfaceSpeed() {
	iSpeed := smart.json.Get("interface_speed")
	for _, speedType := range []string{"max", "current"} {
		tSpeed := iSpeed.Get(speedType)
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceInterfaceSpeed,
			prometheus.GaugeValue,
			tSpeed.Get("units_per_second").Float()*tSpeed.Get("bits_per_unit").Float(),
			smart.device.device,
			speedType,
		)
	}
}

func (smart *SMARTctl) mineDeviceAttribute() {
	for _, attribute := range smart.json.Get("ata_smart_attributes.table").Array() {
		name := strings.TrimSpace(attribute.Get("name").String())
		flagsShort := strings.TrimSpace(attribute.Get("flags.string").String())
		flagsLong := smart.mineLongFlags(attribute.Get("flags"), []string{
			"prefailure",
			"updated_online",
			"performance",
			"error_rate",
			"event_count",
			"auto_keep",
		})
		id := attribute.Get("id").String()
		for key, path := range map[string]string{
			"value":  "value",
			"worst":  "worst",
			"thresh": "thresh",
			"raw":    "raw.value",
		} {
			smart.ch <- prometheus.MustNewConstMetric(
				metricDeviceAttribute,
				prometheus.GaugeValue,
				attribute.Get(path).Float(),
				smart.device.device,
				name,
				flagsShort,
				flagsLong,
				key,
				id,
			)
		}
	}
}

func (smart *SMARTctl) minePowerOnSeconds() {
	pot := smart.json.Get("power_on_time")
	smart.ch <- prometheus.MustNewConstMetric(
		metricDevicePowerOnSeconds,
		prometheus.CounterValue,
		GetFloatIfExists(pot, "hours", 0)*60*60+GetFloatIfExists(pot, "minutes", 0)*60,
		smart.device.device,
	)
}

func (smart *SMARTctl) mineRotationRate() {
	rRate := GetFloatIfExists(smart.json, "rotation_rate", 0)
	if rRate > 0 {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceRotationRate,
			prometheus.GaugeValue,
			rRate,
			smart.device.device,
		)
	}
}

func (smart *SMARTctl) mineTemperatures() {
	temperatures := smart.json.Get("temperature")
	if temperatures.Exists() {
		temperatures.ForEach(func(key, value gjson.Result) bool {
			smart.ch <- prometheus.MustNewConstMetric(
				metricDeviceTemperature,
				prometheus.GaugeValue,
				value.Float(),
				smart.device.device,
				key.String(),
			)
			return true
		})
	}
}

func (smart *SMARTctl) minePowerCycleCount() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDevicePowerCycleCount,
		prometheus.CounterValue,
		smart.json.Get("power_cycle_count").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineDeviceSCTStatus() {
	status := smart.json.Get("ata_sct_status")
	if status.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceState,
			prometheus.GaugeValue,
			status.Get("device_state").Float(),
			smart.device.device,
		)
	}
}

func (smart *SMARTctl) minePercentageUsed() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDevicePercentageUsed,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.percentage_used").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineAvailableSpare() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceAvailableSpare,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.available_spare").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineAvailableSpareThreshold() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceAvailableSpareThreshold,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.available_spare_threshold").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineCriticalWarning() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceCriticalWarning,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.critical_warning").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineMediaErrors() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceMediaErrors,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.media_errors").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineNumErrLogEntries() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceNumErrLogEntries,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.num_err_log_entries").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineBytesRead() {
	blockSize := smart.json.Get("logical_block_size").Float() * 1024
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceBytesRead,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.data_units_read").Float()*blockSize,
		smart.device.device,
	)
}

func (smart *SMARTctl) mineBytesWritten() {
	blockSize := smart.json.Get("logical_block_size").Float() * 1024
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceBytesWritten,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.data_units_written").Float()*blockSize,
		smart.device.device,
	)
}

func (smart *SMARTctl) mineSmartStatus() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceSmartStatus,
		prometheus.GaugeValue,
		smart.json.Get("smart_status.passed").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineDeviceStatistics() {
	for _, page := range smart.json.Get("ata_device_statistics.pages").Array() {
		table := strings.TrimSpace(page.Get("name").String())
		// skip vendor-specific statistics (they lead to duplicate metric labels on Seagate Exos drives,
		// see https://github.com/Sheridan/smartctl_exporter/issues/3 for details)
		if table == "Vendor Specific Statistics" {
			continue
		}
		for _, statistic := range page.Get("table").Array() {
			smart.ch <- prometheus.MustNewConstMetric(
				metricDeviceStatistics,
				prometheus.GaugeValue,
				statistic.Get("value").Float(),
				smart.device.device,
				smart.device.family,
				smart.device.model,
				smart.device.serial,
				table,
				strings.TrimSpace(statistic.Get("name").String()),
				strings.TrimSpace(statistic.Get("flags.string").String()),
				smart.mineLongFlags(statistic.Get("flags"), []string{
					"valid",
					"normalized",
					"supports_dsn",
					"monitored_condition_met",
				}),
			)
		}
	}

	for _, statistic := range smart.json.Get("sata_phy_event_counters.table").Array() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceStatistics,
			prometheus.GaugeValue,
			statistic.Get("value").Float(),
			smart.device.device,
			"SATA PHY Event Counters",
			strings.TrimSpace(statistic.Get("name").String()),
			"V---",
			"valid",
		)
	}
}

func (smart *SMARTctl) mineLongFlags(json gjson.Result, flags []string) string {
	var result []string
	for _, flag := range flags {
		jFlag := json.Get(flag)
		if jFlag.Exists() && jFlag.Bool() {
			result = append(result, flag)
		}
	}
	return strings.Join(result, ",")
}

func (smart *SMARTctl) mineDeviceStatus() {
	status := smart.json.Get("smart_status")
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceStatus,
		prometheus.GaugeValue,
		status.Get("passed").Float(),
		smart.device.device,
	)
}

func (smart *SMARTctl) mineDeviceErrorLog() {
	for logType, status := range smart.json.Get("ata_smart_error_log").Map() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceErrorLogCount,
			prometheus.GaugeValue,
			status.Get("count").Float(),
			smart.device.device,
			logType,
		)
	}
}

func (smart *SMARTctl) mineDeviceSelfTestLog() {
	for logType, status := range smart.json.Get("ata_smart_self_test_log").Map() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceSelfTestLogCount,
			prometheus.GaugeValue,
			status.Get("count").Float(),
			smart.device.device,
			logType,
		)
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceSelfTestLogErrorCount,
			prometheus.GaugeValue,
			status.Get("error_count_total").Float(),
			smart.device.device,
			logType,
		)
	}
}

func (smart *SMARTctl) mineDeviceERC() {
	for ercType, status := range smart.json.Get("ata_sct_erc").Map() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceERCSeconds,
			prometheus.GaugeValue,
			status.Get("deciseconds").Float()/10.0,
			smart.device.device,
			ercType,
		)
	}
}

func (smart *SMARTctl) mineSCSIGrownDefectList() {
	scsi_grown_defect_list := smart.json.Get("scsi_grown_defect_list")
	if scsi_grown_defect_list.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricSCSIGrownDefectList,
			prometheus.GaugeValue,
			scsi_grown_defect_list.Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
		)
	}
}

func (smart *SMARTctl) mineSCSIErrorCounterLog() {
	SCSIHealth := smart.json.Get("scsi_error_counter_log")
	if SCSIHealth.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricReadErrorsCorrectedByRereadsRewrites,
			prometheus.GaugeValue,
			SCSIHealth.Get("read.errors_corrected_by_rereads_rewrites").Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
		)
		smart.ch <- prometheus.MustNewConstMetric(
			metricReadTotalUncorrectedErrors,
			prometheus.GaugeValue,
			SCSIHealth.Get("read.total_uncorrected_errors").Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
		)
		smart.ch <- prometheus.MustNewConstMetric(
			metricWriteErrorsCorrectedByRereadsRewrites,
			prometheus.GaugeValue,
			SCSIHealth.Get("write.errors_corrected_by_rereads_rewrites").Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
		)
		smart.ch <- prometheus.MustNewConstMetric(
			metricWriteTotalUncorrectedErrors,
			prometheus.GaugeValue,
			SCSIHealth.Get("write.total_uncorrected_errors").Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
		)
	}
}
