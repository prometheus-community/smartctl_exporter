## 0.9.2 / 2023-06-10

* [FEATURE] Add device include/exclude filters for the automatic scanning #99
* [FEATURE] Add background scanning for devices #140
* [ENHANCEMENT] Critical metrics for SCSI disks added #131
* [ENHANCEMENT] Added device name to logger rc code parser #141

## 0.9.1 / 2022-11-06

* [BUGFIX] Fix runtime error: index out of range in mineVersion #93
* [BUGFIX] Fix race condition with maps and goroutines #94

## 0.9.0 / 2022-10-20

Breaking Changes:
- Now labels with device model & serial number landed only to smartctl_device meta metric
- /dev/ prefix pruned from device label for matching with node_exporter labels

* [CHANGE] Removed doubled NVMe metrics #82
* [CHANGE] Reduced number of meta labels #83
* [FEATURE] Added disk form_factor meta label #84
* [CHANGE] Pruned /dev/ prefix from device label #88
* [ENHANCEMENT] remove `os.stat` in order to fit in Windows #86
* [ENHANCEMENT] Skip vendor-specific statistics that lead to duplicate metric labels #28

## 0.8.0 / 2022-10-03

Breaking Changes:
All configuration has been moved from the config file to command line flags.

* [CHANGE] Refactor exporter config #68
* [BUGFIX] Fix smartctl command args to avoid wakeups #74
* [ENHANCEMENT] Add smartmontools to container image #51

## 0.7.0 / 2022-08-05

First prometheus-community release.

* [FEATURE] Add various new metrics #14
* [BUGFIX] Return the cached value if it's not time to scan again yet #18
* [BUGFIX] Fix exit code bit parsing #37

## 0.6.0 / 2020-10-29

* Parsing smartctl return code and collect metrics if no errors
* Parsing smartctl messages and collect metrics if no errors

## 0.5.0 / 2019-08-17

* smartctl_device: Device info
* smartctl_device_attribute: Device attributes
* smartctl_device_block_size: Device block size
* smartctl_device_capacity_blocks: Device capacity in blocks
* smartctl_device_capacity_bytes: Device capacity in bytes
* smartctl_device_interface_speed: Device interface speed, bits per second
* smartctl_device_power_cycle_count: Device power cycle count
* smartctl_device_power_on_seconds: Device power on seconds
* smartctl_device_rotation_rate: Device rotation rate
* smartctl_device_smartctl_exit_status: Exit status of smartctl on device
* smartctl_device_statistics: Device statistics
* smartctl_device_temperature: Device temperature celsius
* smartctl_version: smartctl version
