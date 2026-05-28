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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

// mineFarmLog collects Seagate FARM log metrics if present in the JSON.
func (smart *SMARTctl) mineFarmLog() {
	farmLog := smart.json.Get("seagate_farm_log")
	if !farmLog.Exists() {
		return
	}
	smart.logger.Debug("Collecting Seagate FARM log metrics", "device", smart.device.device)
	smart.mineFarmWorkload(farmLog)
	smart.mineFarmErrors(farmLog)
	smart.mineFarmEnvironment(farmLog)
	smart.mineFarmReliability(farmLog)
}

func (smart *SMARTctl) mineFarmWorkload(farmLog gjson.Result) {
	workload := farmLog.Get("page_2_workload_statistics")
	if !workload.Exists() {
		return
	}

	if v := workload.Get("total_read_commands"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmWorkloadReadCommands,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := workload.Get("total_write_commands"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmWorkloadWriteCommands,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := workload.Get("logical_sectors_read"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmWorkloadLogicalSectorsRead,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := workload.Get("logical_sectors_written"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmWorkloadLogicalSectorsWritten,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
}

func (smart *SMARTctl) mineFarmErrors(farmLog gjson.Result) {
	errors := farmLog.Get("page_3_error_statistics")
	if !errors.Exists() {
		return
	}

	if v := errors.Get("number_of_unrecoverable_read_errors"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmErrorUnrecoverableRead,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := errors.Get("number_of_unrecoverable_write_errors"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmErrorUnrecoverableWrite,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := errors.Get("number_of_reallocated_sectors"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmErrorReallocatedSectors,
			prometheus.GaugeValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := errors.Get("number_of_reallocated_candidate_sectors"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmErrorReallocatedCandidates,
			prometheus.GaugeValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := errors.Get("total_crc_errors"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmErrorCRCErrors,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := errors.Get("command_time_out_count_total"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmErrorCommandTimeouts,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}

	// Per-head unrecoverable error metrics
	smart.mineFarmErrorsByHead(errors)
}

func (smart *SMARTctl) mineFarmErrorsByHead(errors gjson.Result) {
	errors.ForEach(func(key, value gjson.Result) bool {
		keyStr := key.String()
		if !strings.HasPrefix(keyStr, "cum_lifetime_unrecoverable_by_head_") {
			return true
		}
		head := strings.TrimPrefix(keyStr, "cum_lifetime_unrecoverable_by_head_")
		if !value.IsObject() {
			return true
		}
		if v := value.Get("cum_lifetime_unrecoverable_read_repeating"); v.Exists() {
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmErrorUnrecoverableByHead,
				prometheus.CounterValue,
				v.Float(),
				smart.device.device,
				head,
				"read_repeating",
			)
		}
		if v := value.Get("cum_lifetime_unrecoverable_read_unique"); v.Exists() {
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmErrorUnrecoverableByHead,
				prometheus.CounterValue,
				v.Float(),
				smart.device.device,
				head,
				"read_unique",
			)
		}
		return true
	})
}

func (smart *SMARTctl) mineFarmEnvironment(farmLog gjson.Result) {
	env := farmLog.Get("page_4_environment_statistics")
	if !env.Exists() {
		return
	}

	tempFields := map[string]string{
		"curent_temp":      "current",
		"highest_temp":     "highest",
		"lowest_temp":      "lowest",
		"average_temp":     "average",
		"average_long_temp": "average_long",
		"highest_short_temp": "highest_short",
		"lowest_short_temp":  "lowest_short",
		"highest_long_temp":  "highest_long",
		"lowest_long_temp":   "lowest_long",
	}
	for field, label := range tempFields {
		if v := env.Get(field); v.Exists() {
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmEnvironmentTemperature,
				prometheus.GaugeValue,
				v.Float(),
				smart.device.device,
				label,
			)
		}
	}

	if v := env.Get("humidity"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmEnvironmentHumidity,
			prometheus.GaugeValue,
			v.Float(),
			smart.device.device,
		)
	}

	if v := env.Get("current_motor_power"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmEnvironmentMotorPower,
			prometheus.GaugeValue,
			v.Float(),
			smart.device.device,
		)
	}

	voltageFields12V := map[string]string{
		"current_12v_in_mv": "current",
		"minimum_12v_in_mv": "minimum",
		"maximum_12v_in_mv": "maximum",
	}
	for field, label := range voltageFields12V {
		if v := env.Get(field); v.Exists() {
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmEnvironment12V,
				prometheus.GaugeValue,
				v.Float(),
				smart.device.device,
				label,
			)
		}
	}

	voltageFields5V := map[string]string{
		"current_5v_in_mv": "current",
		"minimum_5v_in_mv": "minimum",
		"maximum_5v_in_mv": "maximum",
	}
	for field, label := range voltageFields5V {
		if v := env.Get(field); v.Exists() {
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmEnvironment5V,
				prometheus.GaugeValue,
				v.Float(),
				smart.device.device,
				label,
			)
		}
	}
}

func (smart *SMARTctl) mineFarmReliability(farmLog gjson.Result) {
	rel := farmLog.Get("page_5_reliability_statistics")
	if !rel.Exists() {
		return
	}

	if v := rel.Get("attr_error_rate_raw"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmReliabilityErrorRate,
			prometheus.GaugeValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := rel.Get("attr_seek_error_rate_raw"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmReliabilitySeekErrorRate,
			prometheus.GaugeValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := rel.Get("high_priority_unload_events"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmReliabilityHighPriorityUnloads,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}
	if v := rel.Get("helium_presure_trip"); v.Exists() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricFarmReliabilityHeliumPressureTrip,
			prometheus.CounterValue,
			v.Float(),
			smart.device.device,
		)
	}

	// Per-head reliability metrics
	smart.mineFarmReliabilityByHead(rel)
}

func (smart *SMARTctl) mineFarmReliabilityByHead(rel gjson.Result) {
	// Reallocated sectors by head
	rel.ForEach(func(key, value gjson.Result) bool {
		keyStr := key.String()

		if strings.HasPrefix(keyStr, "number_of_reallocated_sectors_by_head_") {
			head := strings.TrimPrefix(keyStr, "number_of_reallocated_sectors_by_head_")
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmReliabilityReallocatedByHead,
				prometheus.GaugeValue,
				value.Float(),
				smart.device.device,
				head,
			)
		} else if strings.HasPrefix(keyStr, "number_of_reallocation_candidate_sectors_by_head_") {
			head := strings.TrimPrefix(keyStr, "number_of_reallocation_candidate_sectors_by_head_")
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmReliabilityCandidatesByHead,
				prometheus.GaugeValue,
				value.Float(),
				smart.device.device,
				head,
			)
		} else if strings.HasPrefix(keyStr, "mr_head_resistance_from_head_") {
			head := strings.TrimPrefix(keyStr, "mr_head_resistance_from_head_")
			smart.ch <- prometheus.MustNewConstMetric(
				metricFarmReliabilityMRHeadResistanceByHead,
				prometheus.GaugeValue,
				value.Float(),
				smart.device.device,
				head,
			)
		}

		return true
	})

	// Skip write detect by head — multiple types per head
	headCount := smart.json.Get("seagate_farm_log.page_1_drive_information.number_of_heads").Int()
	for i := int64(0); i < headCount; i++ {
		head := fmt.Sprintf("%d", i)
		skipTypes := map[string]string{
			fmt.Sprintf("dvga_skip_write_detect_by_head_%d", i):                   "dvga",
			fmt.Sprintf("rvga_skip_write_detect_by_head_%d", i):                   "rvga",
			fmt.Sprintf("fvga_skip_write_detect_by_head_%d", i):                   "fvga",
			fmt.Sprintf("skip_write_detect_threshold_exceeded_by_head_%d", i):     "threshold_exceeded",
		}
		for field, detectType := range skipTypes {
			if v := rel.Get(field); v.Exists() {
				smart.ch <- prometheus.MustNewConstMetric(
					metricFarmReliabilitySkipWriteByHead,
					prometheus.GaugeValue,
					v.Float(),
					smart.device.device,
					head,
					detectType,
				)
			}
		}
	}
}
