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
		},
		nil,
	)
	metricDeviceCapacityBlocks = prometheus.NewDesc(
		"smartctl_device_capacity_blocks",
		"Device capacity in blocks",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceCapacityBytes = prometheus.NewDesc(
		"smartctl_device_capacity_bytes",
		"Device capacity in bytes",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceBlockSize = prometheus.NewDesc(
		"smartctl_device_block_size",
		"Device block size",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"blocks_type",
		},
		nil,
	)
	metricDeviceInterfaceSpeed = prometheus.NewDesc(
		"smartctl_device_interface_speed",
		"Device interface speed, bits per second",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"speed_type",
		},
		nil,
	)
	metricDeviceAttribute = prometheus.NewDesc(
		"smartctl_device_attribute",
		"Device attributes",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
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
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceRotationRate = prometheus.NewDesc(
		"smartctl_device_rotation_rate",
		"Device rotation rate",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceTemperature = prometheus.NewDesc(
		"smartctl_device_temperature",
		"Device temperature celsius",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"temperature_type",
		},
		nil,
	)
	metricDevicePowerCycleCount = prometheus.NewDesc(
		"smartctl_device_power_cycle_count",
		"Device power cycle count",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDevicePercentageUsed = prometheus.NewDesc(
		"smartctl_device_percentage_used",
		"Device write percentage used",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceAvailableSpare = prometheus.NewDesc(
		"smartctl_device_available_spare",
		"Normalized percentage (0 to 100%) of the remaining spare capacity available",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceAvailableSpareThreshold = prometheus.NewDesc(
		"smartctl_device_available_spare_threshold",
		"When the Available Spare falls below the threshold indicated in this field, an asynchronous event completion may occur. The value is indicated as a normalized percentage (0 to 100%)",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceCriticalWarning = prometheus.NewDesc(
		"smartctl_device_critical_warning",
		"This field indicates critical warnings for the state of the controller",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceMediaErrors = prometheus.NewDesc(
		"smartctl_device_media_errors",
		"Contains the number of occurrences where the controller detected an unrecovered data integrity error. Errors such as uncorrectable ECC, CRC checksum failure, or LBA tag mismatch are included in this field",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceNumErrLogEntries = prometheus.NewDesc(
		"smartctl_device_num_err_log_entries",
		"Contains the number of Error Information log entries over the life of the controller",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceBytesRead = prometheus.NewDesc(
		"smartctl_device_bytes_read",
		"",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceBytesWritten = prometheus.NewDesc(
		"smartctl_device_bytes_written",
		"",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceSmartStatus = prometheus.NewDesc(
		"smartctl_device_smart_status",
		"General smart status",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceExitStatus = prometheus.NewDesc(
		"smartctl_device_smartctl_exit_status",
		"Exit status of smartctl on device",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceState = prometheus.NewDesc(
		"smartctl_device_state",
		"Device state (0=active, 1=standby, 2=sleep, 3=dst, 4=offline, 5=sct)",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceStatistics = prometheus.NewDesc(
		"smartctl_device_statistics",
		"Device statistics",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"statistic_table",
			"statistic_table_row",
			"statistic_name",
			"statistic_flags_short",
			"statistic_flags_long",
		},
		nil,
	)
	metricCriticalWarning = prometheus.NewDesc(
		"critical_warning",
		"Critical warning counter",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceStatus = prometheus.NewDesc(
		"smartctl_device_status",
		"Device status",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricAvailableSpare = prometheus.NewDesc(
		"available_spare",
		"Available spare",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceErrorLogCount = prometheus.NewDesc(
		"smartctl_device_error_log_count",
		"Device SMART error log count",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"error_log_type",
		},
		nil,
	)
	metricMediaErrors = prometheus.NewDesc(
		"media_errors",
		"Media errors counter",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceSelfTestLogCount = prometheus.NewDesc(
		"smartctl_device_self_test_log_count",
		"Device SMART self test log count",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"self_test_log_type",
		},
		nil,
	)
	metricSmartStatus = prometheus.NewDesc(
		"smart_status",
		"Smart status",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
		},
		nil,
	)
	metricDeviceSelfTestLogErrorCount = prometheus.NewDesc(
		"smartctl_device_self_test_log_error_count",
		"Device SMART self test log error count",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"self_test_log_type",
		},
		nil,
	)
	metricDeviceERCSeconds = prometheus.NewDesc(
		"smartctl_device_erc_seconds",
		"Device SMART Error Recovery Control Seconds",
		[]string{
			"device",
			"model_family",
			"model_name",
			"serial_number",
			"op_type",
		},
		nil,
	)
)
