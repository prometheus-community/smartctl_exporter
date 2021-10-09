package main

import (
	"fmt"
	"strconv"
	"strings"

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
	device SMARTDevice
}

// NewSMARTctl is smartctl constructor
func NewSMARTctl(json gjson.Result, ch chan<- prometheus.Metric) SMARTctl {
	smart := SMARTctl{}
	smart.ch = ch
	smart.json = json
	smart.device = SMARTDevice{
		device: strings.TrimSpace(smart.json.Get("device.name").String()),
		serial: strings.TrimSpace(smart.json.Get("serial_number").String()),
		family: strings.TrimSpace(smart.json.Get("model_family").String()),
		model:  strings.TrimSpace(smart.json.Get("model_name").String()),
	}
	return smart
}

// Collect metrics
func (smart *SMARTctl) Collect() {
	logger.Verbose("Collecting metrics from %s: %s, %s", smart.device.device, smart.device.family, smart.device.model)
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
	smart.mineNvmeSmartHealthInformationLog()
	smart.mineNvmeSmartStatus()
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

}

func (smart *SMARTctl) mineExitStatus() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceExitStatus,
		prometheus.GaugeValue,
		smart.json.Get("smartctl.exit_status").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
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
	)
}

func (smart *SMARTctl) mineCapacity() {
	capacity := smart.json.Get("user_capacity")
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceCapacityBlocks,
		prometheus.GaugeValue,
		capacity.Get("blocks").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceCapacityBytes,
		prometheus.GaugeValue,
		capacity.Get("bytes").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
	for _, blockType := range []string{"logical", "physical"} {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceBlockSize,
			prometheus.GaugeValue,
			smart.json.Get(fmt.Sprintf("%s_block_size", blockType)).Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
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
			smart.device.family,
			smart.device.model,
			smart.device.serial,
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
				smart.device.family,
				smart.device.model,
				smart.device.serial,
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
		smart.device.family,
		smart.device.model,
		smart.device.serial,
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
			smart.device.family,
			smart.device.model,
			smart.device.serial,
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
				smart.device.family,
				smart.device.model,
				smart.device.serial,
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
		smart.device.family,
		smart.device.model,
		smart.device.serial,
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
			smart.device.family,
			smart.device.model,
			smart.device.serial,
		)
	}
}

func (smart *SMARTctl) minePercentageUsed() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDevicePercentageUsed,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.percentage_used").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineAvailableSpare() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceAvailableSpare,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.available_spare").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineAvailableSpareThreshold() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceAvailableSpareThreshold,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.available_spare_threshold").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineCriticalWarning() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceCriticalWarning,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.critical_warning").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineMediaErrors() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceMediaErrors,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.media_errors").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineNumErrLogEntries() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceNumErrLogEntries,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.num_err_log_entries").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineBytesRead() {
	blockSize := smart.json.Get("logical_block_size").Float() * 1024
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceBytesRead,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.data_units_read").Float()*blockSize,
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineBytesWritten() {
	blockSize := smart.json.Get("logical_block_size").Float() * 1024
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceBytesWritten,
		prometheus.CounterValue,
		smart.json.Get("nvme_smart_health_information_log.data_units_written").Float()*blockSize,
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineSmartStatus() {
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceSmartStatus,
		prometheus.GaugeValue,
		smart.json.Get("smart_status.passed").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineDeviceStatistics() {
	for _, page := range smart.json.Get("ata_device_statistics.pages").Array() {
		table := strings.TrimSpace(page.Get("name").String())
		for row, statistic := range page.Get("table").Array() {
			smart.ch <- prometheus.MustNewConstMetric(
				metricDeviceStatistics,
				prometheus.GaugeValue,
				statistic.Get("value").Float(),
				smart.device.device,
				smart.device.family,
				smart.device.model,
				smart.device.serial,
				table,
				strconv.Itoa(row),
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

	for row, statistic := range smart.json.Get("sata_phy_event_counters.table").Array() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceStatistics,
			prometheus.GaugeValue,
			statistic.Get("value").Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
			"SATA PHY Event Counters",
			strconv.Itoa(row),
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

func (smart *SMARTctl) mineNvmeSmartHealthInformationLog() {
	iHealth := smart.json.Get("nvme_smart_health_information_log")
	smart.ch <- prometheus.MustNewConstMetric(
		metricCriticalWarning,
		prometheus.GaugeValue,
		iHealth.Get("critical_warning").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
	smart.ch <- prometheus.MustNewConstMetric(
		metricAvailableSpare,
		prometheus.GaugeValue,
		iHealth.Get("available_spare").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
	smart.ch <- prometheus.MustNewConstMetric(
		metricMediaErrors,
		prometheus.GaugeValue,
		iHealth.Get("media_errors").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineDeviceStatus() {
	status := smart.json.Get("smart_status")
	smart.ch <- prometheus.MustNewConstMetric(
		metricDeviceStatus,
		prometheus.GaugeValue,
		status.Get("passed").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineNvmeSmartStatus() {
	iStatus := smart.json.Get("smart_status")
	smart.ch <- prometheus.MustNewConstMetric(
		metricSmartStatus,
		prometheus.GaugeValue,
		iStatus.Get("passed").Float(),
		smart.device.device,
		smart.device.family,
		smart.device.model,
		smart.device.serial,
	)
}

func (smart *SMARTctl) mineDeviceErrorLog() {
	for logType, status := range smart.json.Get("ata_smart_error_log").Map() {
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceErrorLogCount,
			prometheus.GaugeValue,
			status.Get("count").Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
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
			smart.device.family,
			smart.device.model,
			smart.device.serial,
			logType,
		)
		smart.ch <- prometheus.MustNewConstMetric(
			metricDeviceSelfTestLogErrorCount,
			prometheus.GaugeValue,
			status.Get("error_count_total").Float(),
			smart.device.device,
			smart.device.family,
			smart.device.model,
			smart.device.serial,
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
			smart.device.family,
			smart.device.model,
			smart.device.serial,
			ercType,
		)
	}
}
