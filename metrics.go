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
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricSmartctlVersion = prometheus.NewDesc(
		"smartctl_version",
		"smartctl version",
		[]string{
			"json_format_version",
			"smartctl_version",
			"svn_revision",
			"build_info",
		},
		nil,
	)
	metricDeviceModel = prometheus.NewDesc(
		"smartctl_device",
		"Device info",
		[]string{
			"device",
			"interface",
			"protocol",
			"model_family",
			"model_name",
			"serial_number",
			"ata_additional_product_id",
			"firmware_version",
			"ata_version",
			"sata_version",
			"form_factor",
			// scsi_model_name is mapped into model_name
			"scsi_vendor",
			"scsi_product",
			"scsi_revision",
			"scsi_version",
		},
		nil,
	)
	metricDeviceCount = prometheus.NewDesc(
		"smartctl_devices",
		"Number of devices configured or dynamically discovered",
		[]string{},
		nil,
	)
	metricDeviceCapacityBlocks = prometheus.NewDesc(
		"smartctl_device_capacity_blocks",
		"Device capacity in blocks",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceCapacityBytes = prometheus.NewDesc(
		"smartctl_device_capacity_bytes",
		"Device capacity in bytes",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceTotalCapacityBytes = prometheus.NewDesc(
		"smartctl_device_nvme_capacity_bytes",
		"NVMe device total capacity bytes",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceBlockSize = prometheus.NewDesc(
		"smartctl_device_block_size",
		"Device block size",
		[]string{
			"device",
			"blocks_type",
		},
		nil,
	)
	metricDeviceInterfaceSpeed = prometheus.NewDesc(
		"smartctl_device_interface_speed",
		"Device interface speed, bits per second",
		[]string{
			"device",
			"speed_type",
		},
		nil,
	)
	metricDeviceAttribute = prometheus.NewDesc(
		"smartctl_device_attribute",
		"Device attributes",
		[]string{
			"device",
			"attribute_name",
			"attribute_flags_short",
			"attribute_flags_long",
			"attribute_value_type",
			"attribute_id",
		},
		nil,
	)
	metricDevicePowerOnSeconds = prometheus.NewDesc(
		"smartctl_device_power_on_seconds",
		"Device power on seconds",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceRotationRate = prometheus.NewDesc(
		"smartctl_device_rotation_rate",
		"Device rotation rate",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceTemperature = prometheus.NewDesc(
		"smartctl_device_temperature",
		"Device temperature celsius",
		[]string{
			"device",
			"temperature_type",
		},
		nil,
	)
	metricDevicePowerCycleCount = prometheus.NewDesc(
		"smartctl_device_power_cycle_count",
		"Device power cycle count",
		[]string{
			"device",
		},
		nil,
	)
	metricDevicePercentageUsed = prometheus.NewDesc(
		"smartctl_device_percentage_used",
		"Device write percentage used",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceAvailableSpare = prometheus.NewDesc(
		"smartctl_device_available_spare",
		"Normalized percentage (0 to 100%) of the remaining spare capacity available",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceAvailableSpareThreshold = prometheus.NewDesc(
		"smartctl_device_available_spare_threshold",
		"When the Available Spare falls below the threshold indicated in this field, an asynchronous event completion may occur. The value is indicated as a normalized percentage (0 to 100%)",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceCriticalWarning = prometheus.NewDesc(
		"smartctl_device_critical_warning",
		"This field indicates critical warnings for the state of the controller",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceMediaErrors = prometheus.NewDesc(
		"smartctl_device_media_errors",
		"Contains the number of occurrences where the controller detected an unrecovered data integrity error. Errors such as uncorrectable ECC, CRC checksum failure, or LBA tag mismatch are included in this field",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceNumErrLogEntries = prometheus.NewDesc(
		"smartctl_device_num_err_log_entries",
		"Contains the number of Error Information log entries over the life of the controller",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceBytesRead = prometheus.NewDesc(
		"smartctl_device_bytes_read",
		"",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceBytesWritten = prometheus.NewDesc(
		"smartctl_device_bytes_written",
		"",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceSmartStatus = prometheus.NewDesc(
		"smartctl_device_smart_status",
		"General smart status",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceExitStatus = prometheus.NewDesc(
		"smartctl_device_smartctl_exit_status",
		"Exit status of smartctl on device",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceState = prometheus.NewDesc(
		"smartctl_device_state",
		"Device state (0=active, 1=standby, 2=sleep, 3=dst, 4=offline, 5=sct)",
		[]string{
			"device",
		},
		nil,
	)
	metricDeviceStatistics = prometheus.NewDesc(
		"smartctl_device_statistics",
		"Device statistics",
		[]string{
			"device",
			"statistic_table",
			"statistic_name",
			"statistic_flags_short",
			"statistic_flags_long",
		},
		nil,
	)
	metricDeviceErrorLogCount = prometheus.NewDesc(
		"smartctl_device_error_log_count",
		"Device SMART error log count",
		[]string{
			"device",
			"error_log_type",
		},
		nil,
	)
	metricDeviceSelfTestLogCount = prometheus.NewDesc(
		"smartctl_device_self_test_log_count",
		"Device SMART self test log count",
		[]string{
			"device",
			"self_test_log_type",
		},
		nil,
	)
	metricDeviceSelfTestLogErrorCount = prometheus.NewDesc(
		"smartctl_device_self_test_log_error_count",
		"Device SMART self test log error count",
		[]string{
			"device",
			"self_test_log_type",
		},
		nil,
	)
	metricDeviceERCSeconds = prometheus.NewDesc(
		"smartctl_device_erc_seconds",
		"Device SMART Error Recovery Control Seconds",
		[]string{
			"device",
			"op_type",
		},
		nil,
	)
	metricSCSIGrownDefectList = prometheus.NewDesc(
		"smartctl_scsi_grown_defect_list",
		"Device SCSI grown defect list counter",
		[]string{
			"device",
		},
		nil,
	)
	metricReadErrorsCorrectedByRereadsRewrites = prometheus.NewDesc(
		"smartctl_read_errors_corrected_by_rereads_rewrites",
		"Read Errors Corrected by ReReads/ReWrites",
		[]string{
			"device",
		},
		nil,
	)
	metricReadErrorsCorrectedByEccFast = prometheus.NewDesc(
		"smartctl_read_errors_corrected_by_eccfast",
		"Read Errors Corrected by ECC Fast",
		[]string{
			"device",
		},
		nil,
	)
	metricReadErrorsCorrectedByEccDelayed = prometheus.NewDesc(
		"smartctl_read_errors_corrected_by_eccdelayed",
		"Read Errors Corrected by ECC Delayed",
		[]string{
			"device",
		},
		nil,
	)
	metricReadTotalUncorrectedErrors = prometheus.NewDesc(
		"smartctl_read_total_uncorrected_errors",
		"Read Total Uncorrected Errors",
		[]string{
			"device",
		},
		nil,
	)
	metricWriteErrorsCorrectedByRereadsRewrites = prometheus.NewDesc(
		"smartctl_write_errors_corrected_by_rereads_rewrites",
		"Write Errors Corrected by ReReads/ReWrites",
		[]string{
			"device",
		},
		nil,
	)
	metricWriteErrorsCorrectedByEccFast = prometheus.NewDesc(
		"smartctl_write_errors_corrected_by_eccfast",
		"Write Errors Corrected by ECC Fast",
		[]string{
			"device",
		},
		nil,
	)
	metricWriteErrorsCorrectedByEccDelayed = prometheus.NewDesc(
		"smartctl_write_errors_corrected_by_eccdelayed",
		"Write Errors Corrected by ECC Delayed",
		[]string{
			"device",
		},
		nil,
	)
	metricWriteTotalUncorrectedErrors = prometheus.NewDesc(
		"smartctl_write_total_uncorrected_errors",
		"Write Total Uncorrected Errors",
		[]string{
			"device",
		},
		nil,
	)

	// Seagate FARM log metrics
	metricFarmWorkloadReadCommands = prometheus.NewDesc(
		"smartctl_device_farm_workload_read_commands_total",
		"Seagate FARM total read commands",
		[]string{"device"},
		nil,
	)
	metricFarmWorkloadWriteCommands = prometheus.NewDesc(
		"smartctl_device_farm_workload_write_commands_total",
		"Seagate FARM total write commands",
		[]string{"device"},
		nil,
	)
	metricFarmWorkloadLogicalSectorsRead = prometheus.NewDesc(
		"smartctl_device_farm_workload_logical_sectors_read_total",
		"Seagate FARM logical sectors read",
		[]string{"device"},
		nil,
	)
	metricFarmWorkloadLogicalSectorsWritten = prometheus.NewDesc(
		"smartctl_device_farm_workload_logical_sectors_written_total",
		"Seagate FARM logical sectors written",
		[]string{"device"},
		nil,
	)
	metricFarmErrorUnrecoverableRead = prometheus.NewDesc(
		"smartctl_device_farm_error_unrecoverable_read_total",
		"Seagate FARM unrecoverable read errors",
		[]string{"device"},
		nil,
	)
	metricFarmErrorUnrecoverableWrite = prometheus.NewDesc(
		"smartctl_device_farm_error_unrecoverable_write_total",
		"Seagate FARM unrecoverable write errors",
		[]string{"device"},
		nil,
	)
	metricFarmErrorReallocatedSectors = prometheus.NewDesc(
		"smartctl_device_farm_error_reallocated_sectors",
		"Seagate FARM reallocated sectors",
		[]string{"device"},
		nil,
	)
	metricFarmErrorReallocatedCandidates = prometheus.NewDesc(
		"smartctl_device_farm_error_reallocated_candidate_sectors",
		"Seagate FARM reallocated candidate sectors",
		[]string{"device"},
		nil,
	)
	metricFarmErrorCRCErrors = prometheus.NewDesc(
		"smartctl_device_farm_error_crc_errors_total",
		"Seagate FARM total CRC errors",
		[]string{"device"},
		nil,
	)
	metricFarmErrorCommandTimeouts = prometheus.NewDesc(
		"smartctl_device_farm_error_command_timeouts_total",
		"Seagate FARM command timeout count",
		[]string{"device"},
		nil,
	)
	metricFarmEnvironmentTemperature = prometheus.NewDesc(
		"smartctl_device_farm_environment_temperature_celsius",
		"Seagate FARM environment temperature",
		[]string{"device", "temperature_type"},
		nil,
	)
	metricFarmEnvironmentHumidity = prometheus.NewDesc(
		"smartctl_device_farm_environment_humidity_percent",
		"Seagate FARM environment humidity percentage",
		[]string{"device"},
		nil,
	)
	metricFarmEnvironment12V = prometheus.NewDesc(
		"smartctl_device_farm_environment_12v_millivolts",
		"Seagate FARM 12V rail millivolts",
		[]string{"device", "voltage_type"},
		nil,
	)
	metricFarmEnvironment5V = prometheus.NewDesc(
		"smartctl_device_farm_environment_5v_millivolts",
		"Seagate FARM 5V rail millivolts",
		[]string{"device", "voltage_type"},
		nil,
	)
	metricFarmEnvironmentMotorPower = prometheus.NewDesc(
		"smartctl_device_farm_environment_motor_power",
		"Seagate FARM current motor power",
		[]string{"device"},
		nil,
	)
	metricFarmReliabilityErrorRate = prometheus.NewDesc(
		"smartctl_device_farm_reliability_error_rate_raw",
		"Seagate FARM error rate raw value",
		[]string{"device"},
		nil,
	)
	metricFarmReliabilitySeekErrorRate = prometheus.NewDesc(
		"smartctl_device_farm_reliability_seek_error_rate_raw",
		"Seagate FARM seek error rate raw value",
		[]string{"device"},
		nil,
	)
	metricFarmReliabilityHighPriorityUnloads = prometheus.NewDesc(
		"smartctl_device_farm_reliability_high_priority_unload_events",
		"Seagate FARM high priority unload events",
		[]string{"device"},
		nil,
	)
	metricFarmReliabilityHeliumPressureTrip = prometheus.NewDesc(
		"smartctl_device_farm_reliability_helium_pressure_trip",
		"Seagate FARM helium pressure trip count",
		[]string{"device"},
		nil,
	)
	metricFarmErrorUnrecoverableByHead = prometheus.NewDesc(
		"smartctl_device_farm_error_unrecoverable_by_head",
		"Seagate FARM per-head unrecoverable errors",
		[]string{"device", "head", "error_type"},
		nil,
	)
	metricFarmReliabilityReallocatedByHead = prometheus.NewDesc(
		"smartctl_device_farm_reliability_reallocated_sectors_by_head",
		"Seagate FARM per-head reallocated sectors",
		[]string{"device", "head"},
		nil,
	)
	metricFarmReliabilityCandidatesByHead = prometheus.NewDesc(
		"smartctl_device_farm_reliability_reallocation_candidates_by_head",
		"Seagate FARM per-head reallocation candidate sectors",
		[]string{"device", "head"},
		nil,
	)
	metricFarmReliabilitySkipWriteByHead = prometheus.NewDesc(
		"smartctl_device_farm_reliability_skip_write_detect_by_head",
		"Seagate FARM per-head skip write detect count",
		[]string{"device", "head", "detect_type"},
		nil,
	)
	metricFarmReliabilityMRHeadResistanceByHead = prometheus.NewDesc(
		"smartctl_device_farm_reliability_mr_head_resistance_by_head",
		"Seagate FARM per-head MR head resistance",
		[]string{"device", "head"},
		nil,
	)
)
