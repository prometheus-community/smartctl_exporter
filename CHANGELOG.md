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
