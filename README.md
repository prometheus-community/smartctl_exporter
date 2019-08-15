# smartctl_exporter
Export smartctl statistics to prometheus

## Need more?
**If you need additional metrics - contact me :)**
**Create a feature request, describe the metric that you would like to have and attach exported from smartctl json file**

# Configuration
## Command line options
* `--config=/path/to/file.yaml`: Path to configuration file, defaulr `/etc/smartctl_exporter.yaml`
* `--verbose`: verbosed log, default no
* `--debug`: Debug logging, default no
* `--version`: Show version and exit

## Configuration file
Example content:
```
smartctl_exporter:
  bind_to: "[::1]:9633"
  url_path: "/metrics"
  fake_json: no
  smartctl_location: /usr/sbin/smartctl
  collect_not_more_than_period: 120s
  devices:
  - /dev/sda
  - /dev/sdb
  - /dev/sdc
  - /dev/sdd
  - /dev/sde
  - /dev/sdf
```
`fake_json` used for debugging.

# Example metrics
```
# HELP smartctl_device Device info
# TYPE smartctl_device gauge
smartctl_device{ata_additional_product_id="HC7020E0",ata_version="ATA8-ACS, ATA/ATAPI-7 T13/1532D revision 4a",device="/dev/sda",firmware_version="1.03",interface="sat",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",protocol="ATA",sata_version="SATA 3.1",serial_number="P02448109994"} 1.0
smartctl_device{ata_additional_product_id="unknown",ata_version="ACS-2 (minor revision not indicated)",device="/dev/sdc",firmware_version="82.00A82",interface="sat",model_family="Western Digital Red",model_name="WDC WD20EFRX-68EUZN0",protocol="ATA",sata_version="SATA 3.0",serial_number="WD-WCC4M4VX3C69"} 1.0
# HELP smartctl_device_attribute Device attributes
# TYPE smartctl_device_attribute gauge
smartctl_device_attribute{device="/dev/sda",flags="-O----",id="9",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",name="Power_On_Hours",serial_number="P02448109994",value_type="raw"} 2952.0
smartctl_device_attribute{device="/dev/sda",flags="-O----",id="9",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",name="Power_On_Hours",serial_number="P02448109994",value_type="thresh"} 0.0
smartctl_device_attribute{device="/dev/sda",flags="-O----",id="9",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",name="Power_On_Hours",serial_number="P02448109994",value_type="value"} 100.0
smartctl_device_attribute{device="/dev/sda",flags="-O----",id="9",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",name="Power_On_Hours",serial_number="P02448109994",value_type="worst"} 100.0
# HELP smartctl_device_block_size Device block size
# TYPE smartctl_device_block_size gauge
smartctl_device_block_size{blocks_type="logical",device="/dev/sda",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",serial_number="P02448109994"} 512.0
smartctl_device_block_size{blocks_type="logical",device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361"} 512.0
smartctl_device_block_size{blocks_type="logical",device="/dev/sdc",model_family="Western Digital Red",model_name="WDC WD20EFRX-68EUZN0",serial_number="WD-WCC4M4VX3C69"} 512.0
smartctl_device_block_size{blocks_type="logical",device="/dev/sdd",model_family="Western Digital Red",model_name="WDC WD20EFRX-68EUZN0",serial_number="WD-WCC4M6DCPPC7"} 512.0
# HELP smartctl_device_capacity_blocks Device capacity in blocks
# TYPE smartctl_device_capacity_blocks gauge
smartctl_device_capacity_blocks{device="/dev/sda",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",serial_number="P02448109994"} 2.5006968e+08
smartctl_device_capacity_blocks{device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361"} 2.34441648e+08
smartctl_device_capacity_blocks{device="/dev/sdc",model_family="Western Digital Red",model_name="WDC WD20EFRX-68EUZN0",serial_number="WD-WCC4M4VX3C69"} 3.907029168e+09
# HELP smartctl_device_capacity_bytes Device capacity in bytes
# TYPE smartctl_device_capacity_bytes gauge
smartctl_device_capacity_bytes{device="/dev/sda",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",serial_number="P02448109994"} 1.2803567616e+11
smartctl_device_capacity_bytes{device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361"} 1.20034123776e+11
smartctl_device_capacity_bytes{device="/dev/sdc",model_family="Western Digital Red",model_name="WDC WD20EFRX-68EUZN0",serial_number="WD-WCC4M4VX3C69"} 2.000398934016e+12
# HELP smartctl_device_interface_speed Device interface speed, bits per second
# TYPE smartctl_device_interface_speed gauge
smartctl_device_interface_speed{device="/dev/sda",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",serial_number="P02448109994",speed_type="current"} 6e+09
smartctl_device_interface_speed{device="/dev/sda",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",serial_number="P02448109994",speed_type="max"} 6e+09
smartctl_device_interface_speed{device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361",speed_type="current"} 6e+09
# HELP smartctl_device_power_on_seconds Device power on seconds
# TYPE smartctl_device_power_on_seconds counter
smartctl_device_power_on_seconds{device="/dev/sda",model_family="Plextor M3/M5/M6 Series SSDs",model_name="PLEXTOR PX-128M6S",serial_number="P02448109994"} 1.06272e+07
smartctl_device_power_on_seconds{device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361"} 1.0568592e+08
smartctl_device_power_on_seconds{device="/dev/sdc",model_family="Western Digital Red",model_name="WDC WD20EFRX-68EUZN0",serial_number="WD-WCC4M4VX3C69"} 6.68232e+07
# HELP smartctl_device_temperature Device temperature celsius
# TYPE smartctl_device_temperature gauge
smartctl_device_temperature{device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361",temperature_type="current"} 30.0
smartctl_device_temperature{device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361",temperature_type="lifetime_max"} 30.0
smartctl_device_temperature{device="/dev/sdb",model_family="SandForce Driven SSDs",model_name="OCZ-VERTEX3",serial_number="A20Y8011312000361",temperature_type="lifetime_min"} 30.0
```
