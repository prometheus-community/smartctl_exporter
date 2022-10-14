# Example output
```
# HELP smartctl_device Device info
# TYPE smartctl_device gauge
smartctl_device{ata_additional_product_id="unknown",ata_version="",device="/dev/nvme0",firmware_version="1.03",form_factor="",interface="nvme",model_family="",model_name="PLEXTOR PX-256M9PY +",protocol="NVMe",sata_version="",serial_number="P02050305257"} 1
smartctl_device{ata_additional_product_id="unknown",ata_version="ACS-2, ACS-3 T13/2161-D revision 3b",device="/dev/sdd",firmware_version="1.14",form_factor="",interface="sat",model_family="",model_name="ADATA SP920SS",protocol="ATA",sata_version="SATA 3.0",serial_number="8H0920011311"} 1
smartctl_device{ata_additional_product_id="unknown",ata_version="ACS-3 T13/2161-D revision 5",device="/dev/sda",firmware_version="0103",form_factor="3.5 inches",interface="sat",model_family="Toshiba MG07ACA... Enterprise Capacity HDD",model_name="TOSHIBA MG07ACA12TE",protocol="ATA",sata_version="SATA 3.3",serial_number="X010A0GTF96G"} 1
smartctl_device{ata_additional_product_id="unknown",ata_version="ACS-3 T13/2161-D revision 5",device="/dev/sdb",firmware_version="0103",form_factor="3.5 inches",interface="sat",model_family="Toshiba MG07ACA... Enterprise Capacity HDD",model_name="TOSHIBA MG07ACA12TE",protocol="ATA",sata_version="SATA 3.3",serial_number="X010A069F96G"} 1
smartctl_device{ata_additional_product_id="unknown",ata_version="ACS-3 T13/2161-D revision 5",device="/dev/sdc",firmware_version="0103",form_factor="3.5 inches",interface="sat",model_family="Toshiba MG07ACA... Enterprise Capacity HDD",model_name="TOSHIBA MG07ACA12TE",protocol="ATA",sata_version="SATA 3.3",serial_number="X010A0QRF96G"} 1
# HELP smartctl_device_attribute Device attributes
# TYPE smartctl_device_attribute gauge
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="177",attribute_name="Wear_Leveling_Count",attribute_value_type="raw",device="/dev/sdd"} 8
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="177",attribute_name="Wear_Leveling_Count",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="177",attribute_name="Wear_Leveling_Count",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="177",attribute_name="Wear_Leveling_Count",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="233",attribute_name="Media_Wearout_Indicator",attribute_value_type="raw",device="/dev/sdd"} 1.3907933e+07
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="233",attribute_name="Media_Wearout_Indicator",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="233",attribute_name="Media_Wearout_Indicator",attribute_value_type="value",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="",attribute_flags_short="------",attribute_id="233",attribute_name="Media_Wearout_Indicator",attribute_value_type="worst",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="174",attribute_name="Unknown_Attribute",attribute_value_type="raw",device="/dev/sdd"} 29
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="174",attribute_name="Unknown_Attribute",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="174",attribute_name="Unknown_Attribute",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="174",attribute_name="Unknown_Attribute",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="event_count,auto_keep",attribute_flags_short="----CK",attribute_id="198",attribute_name="Offline_Uncorrectable",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="performance,error_rate,event_count",attribute_flags_short="--SRC-",attribute_id="195",attribute_name="Hardware_ECC_Recovered",attribute_value_type="raw",device="/dev/sdd"} 12817
smartctl_device_attribute{attribute_flags_long="performance,error_rate,event_count",attribute_flags_short="--SRC-",attribute_id="195",attribute_name="Hardware_ECC_Recovered",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="performance,error_rate,event_count",attribute_flags_short="--SRC-",attribute_id="195",attribute_name="Hardware_ECC_Recovered",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="performance,error_rate,event_count",attribute_flags_short="--SRC-",attribute_id="195",attribute_name="Hardware_ECC_Recovered",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="thresh",device="/dev/sda"} 1
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="thresh",device="/dev/sdb"} 1
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="thresh",device="/dev/sdc"} 1
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure",attribute_flags_short="P-----",attribute_id="240",attribute_name="Head_Flying_Hours",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="thresh",device="/dev/sda"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="thresh",device="/dev/sdb"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="thresh",device="/dev/sdc"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="2",attribute_name="Throughput_Performance",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="thresh",device="/dev/sda"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="thresh",device="/dev/sdb"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="thresh",device="/dev/sdc"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,performance",attribute_flags_short="P-S---",attribute_id="8",attribute_name="Seek_Time_Performance",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="thresh",device="/dev/sda"} 75
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="thresh",device="/dev/sdb"} 75
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="thresh",device="/dev/sdc"} 75
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="23",attribute_name="Helium_Condition_Lower",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="thresh",device="/dev/sda"} 75
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="thresh",device="/dev/sdb"} 75
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="thresh",device="/dev/sdc"} 75
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,auto_keep",attribute_flags_short="PO---K",attribute_id="24",attribute_name="Helium_Condition_Upper",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="thresh",device="/dev/sda"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="thresh",device="/dev/sdb"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="thresh",device="/dev/sdc"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="thresh",device="/dev/sda"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="thresh",device="/dev/sdb"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="thresh",device="/dev/sdc"} 50
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,error_rate",attribute_flags_short="PO-R--",attribute_id="7",attribute_name="Seek_Error_Rate",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count",attribute_flags_short="PO--C-",attribute_id="231",attribute_name="Unknown_SSD_Attribute",attribute_value_type="raw",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count",attribute_flags_short="PO--C-",attribute_id="231",attribute_name="Unknown_SSD_Attribute",attribute_value_type="thresh",device="/dev/sdd"} 10
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count",attribute_flags_short="PO--C-",attribute_id="231",attribute_name="Unknown_SSD_Attribute",attribute_value_type="value",device="/dev/sdd"} 97
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count",attribute_flags_short="PO--C-",attribute_id="231",attribute_name="Unknown_SSD_Attribute",attribute_value_type="worst",device="/dev/sdd"} 97
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="thresh",device="/dev/sda"} 30
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="thresh",device="/dev/sdb"} 30
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="thresh",device="/dev/sdc"} 30
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="10",attribute_name="Spin_Retry_Count",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="raw",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="thresh",device="/dev/sda"} 10
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="thresh",device="/dev/sdb"} 10
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="thresh",device="/dev/sdc"} 10
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="thresh",device="/dev/sdd"} 3
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,event_count,auto_keep",attribute_flags_short="PO--CK",attribute_id="5",attribute_name="Reallocated_Sector_Ct",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="raw",device="/dev/sda"} 3977
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="raw",device="/dev/sdb"} 4001
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="raw",device="/dev/sdc"} 4040
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="thresh",device="/dev/sda"} 1
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="thresh",device="/dev/sdb"} 1
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="thresh",device="/dev/sdc"} 1
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,auto_keep",attribute_flags_short="POS--K",attribute_id="3",attribute_name="Spin_Up_Time",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,error_rate",attribute_flags_short="POSR--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="raw",device="/dev/sdd"} 1352
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,error_rate",attribute_flags_short="POSR--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,error_rate",attribute_flags_short="POSR--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="prefailure,updated_online,performance,error_rate",attribute_flags_short="POSR--",attribute_id="1",attribute_name="Raw_Read_Error_Rate",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="raw",device="/dev/sda"} 1.835008e+06
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="raw",device="/dev/sdb"} 1.31072e+06
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="raw",device="/dev/sdc"} 262145
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online",attribute_flags_short="-O----",attribute_id="220",attribute_name="Disk_Shift",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="raw",device="/dev/sda"} 1.8468503556e+11
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="raw",device="/dev/sdb"} 1.8468503556e+11
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="raw",device="/dev/sdc"} 1.84684773415e+11
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="raw",device="/dev/sdd"} 3.342375e+06
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="value",device="/dev/sdd"} 39
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="194",attribute_name="Temperature_Celsius",attribute_value_type="worst",device="/dev/sdd"} 51
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,auto_keep",attribute_flags_short="-O---K",attribute_id="224",attribute_name="Load_Friction",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="raw",device="/dev/sda"} 11
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="raw",device="/dev/sdb"} 11
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="raw",device="/dev/sdc"} 11
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="raw",device="/dev/sdd"} 44
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="12",attribute_name="Power_Cycle_Count",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="187",attribute_name="Reported_Uncorrect",attribute_value_type="raw",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="187",attribute_name="Reported_Uncorrect",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="187",attribute_name="Reported_Uncorrect",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="187",attribute_name="Reported_Uncorrect",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="raw",device="/dev/sdb"} 1
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="raw",device="/dev/sdc"} 1
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="191",attribute_name="G-Sense_Error_Rate",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="raw",device="/dev/sda"} 6
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="raw",device="/dev/sdb"} 6
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="raw",device="/dev/sdc"} 6
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="192",attribute_name="Power-Off_Retract_Count",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="raw",device="/dev/sda"} 125
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="raw",device="/dev/sdb"} 151
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="raw",device="/dev/sdc"} 158
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="193",attribute_name="Load_Cycle_Count",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="196",attribute_name="Reallocated_Event_Count",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="197",attribute_name="Current_Pending_Sector",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="value",device="/dev/sda"} 200
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="value",device="/dev/sdb"} 200
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="value",device="/dev/sdc"} 200
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="worst",device="/dev/sda"} 200
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="worst",device="/dev/sdb"} 200
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="199",attribute_name="UDMA_CRC_Error_Count",attribute_value_type="worst",device="/dev/sdc"} 200
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="raw",device="/dev/sda"} 13299
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="raw",device="/dev/sdb"} 13348
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="raw",device="/dev/sdc"} 13356
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="value",device="/dev/sda"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="value",device="/dev/sdb"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="value",device="/dev/sdc"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="worst",device="/dev/sda"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="worst",device="/dev/sdb"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="222",attribute_name="Loaded_Hours",attribute_value_type="worst",device="/dev/sdc"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="raw",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="raw",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="raw",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="223",attribute_name="Load_Retry_Count",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="234",attribute_name="Unknown_Attribute",attribute_value_type="raw",device="/dev/sdd"} 3.034552e+06
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="234",attribute_name="Unknown_Attribute",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="234",attribute_name="Unknown_Attribute",attribute_value_type="value",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="234",attribute_name="Unknown_Attribute",attribute_value_type="worst",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="241",attribute_name="Total_LBAs_Written",attribute_value_type="raw",device="/dev/sdd"} 5.101427e+06
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="241",attribute_name="Total_LBAs_Written",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="241",attribute_name="Total_LBAs_Written",attribute_value_type="value",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="241",attribute_name="Total_LBAs_Written",attribute_value_type="worst",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="242",attribute_name="Total_LBAs_Read",attribute_value_type="raw",device="/dev/sdd"} 4.84436e+06
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="242",attribute_name="Total_LBAs_Read",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="242",attribute_name="Total_LBAs_Read",attribute_value_type="value",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="242",attribute_name="Total_LBAs_Read",attribute_value_type="worst",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="raw",device="/dev/sda"} 11
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="raw",device="/dev/sdb"} 11
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="raw",device="/dev/sdc"} 11
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="4",attribute_name="Start_Stop_Count",attribute_value_type="worst",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="raw",device="/dev/sda"} 13359
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="raw",device="/dev/sdb"} 13407
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="raw",device="/dev/sdc"} 13432
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="raw",device="/dev/sdd"} 37412
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="thresh",device="/dev/sdd"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="value",device="/dev/sda"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="value",device="/dev/sdb"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="value",device="/dev/sdc"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="value",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="worst",device="/dev/sda"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="worst",device="/dev/sdb"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="worst",device="/dev/sdc"} 67
smartctl_device_attribute{attribute_flags_long="updated_online,event_count,auto_keep",attribute_flags_short="-O--CK",attribute_id="9",attribute_name="Power_On_Hours",attribute_value_type="worst",device="/dev/sdd"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="raw",device="/dev/sda"} 536
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="raw",device="/dev/sdb"} 532
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="raw",device="/dev/sdc"} 533
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="thresh",device="/dev/sda"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="thresh",device="/dev/sdb"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="thresh",device="/dev/sdc"} 0
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="value",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="value",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="value",device="/dev/sdc"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="worst",device="/dev/sda"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="worst",device="/dev/sdb"} 100
smartctl_device_attribute{attribute_flags_long="updated_online,performance,auto_keep",attribute_flags_short="-OS--K",attribute_id="226",attribute_name="Load-in_Time",attribute_value_type="worst",device="/dev/sdc"} 100
# HELP smartctl_device_available_spare Normalized percentage (0 to 100%) of the remaining spare capacity available
# TYPE smartctl_device_available_spare counter
smartctl_device_available_spare{device="/dev/nvme0"} 100
smartctl_device_available_spare{device="/dev/sda"} 0
smartctl_device_available_spare{device="/dev/sdb"} 0
smartctl_device_available_spare{device="/dev/sdc"} 0
smartctl_device_available_spare{device="/dev/sdd"} 0
# HELP smartctl_device_available_spare_threshold When the Available Spare falls below the threshold indicated in this field, an asynchronous event completion may occur. The value is indicated as a normalized percentage (0 to 100%)
# TYPE smartctl_device_available_spare_threshold counter
smartctl_device_available_spare_threshold{device="/dev/nvme0"} 0
smartctl_device_available_spare_threshold{device="/dev/sda"} 0
smartctl_device_available_spare_threshold{device="/dev/sdb"} 0
smartctl_device_available_spare_threshold{device="/dev/sdc"} 0
smartctl_device_available_spare_threshold{device="/dev/sdd"} 0
# HELP smartctl_device_block_size Device block size
# TYPE smartctl_device_block_size gauge
smartctl_device_block_size{blocks_type="logical",device="/dev/nvme0"} 512
smartctl_device_block_size{blocks_type="logical",device="/dev/sda"} 512
smartctl_device_block_size{blocks_type="logical",device="/dev/sdb"} 512
smartctl_device_block_size{blocks_type="logical",device="/dev/sdc"} 512
smartctl_device_block_size{blocks_type="logical",device="/dev/sdd"} 512
smartctl_device_block_size{blocks_type="physical",device="/dev/nvme0"} 0
smartctl_device_block_size{blocks_type="physical",device="/dev/sda"} 4096
smartctl_device_block_size{blocks_type="physical",device="/dev/sdb"} 4096
smartctl_device_block_size{blocks_type="physical",device="/dev/sdc"} 4096
smartctl_device_block_size{blocks_type="physical",device="/dev/sdd"} 512
# HELP smartctl_device_bytes_read 
# TYPE smartctl_device_bytes_read counter
smartctl_device_bytes_read{device="/dev/nvme0"} 5.1499918426112e+13
smartctl_device_bytes_read{device="/dev/sda"} 0
smartctl_device_bytes_read{device="/dev/sdb"} 0
smartctl_device_bytes_read{device="/dev/sdc"} 0
smartctl_device_bytes_read{device="/dev/sdd"} 0
# HELP smartctl_device_bytes_written 
# TYPE smartctl_device_bytes_written counter
smartctl_device_bytes_written{device="/dev/nvme0"} 1.9632806690816e+13
smartctl_device_bytes_written{device="/dev/sda"} 0
smartctl_device_bytes_written{device="/dev/sdb"} 0
smartctl_device_bytes_written{device="/dev/sdc"} 0
smartctl_device_bytes_written{device="/dev/sdd"} 0
# HELP smartctl_device_capacity_blocks Device capacity in blocks
# TYPE smartctl_device_capacity_blocks gauge
smartctl_device_capacity_blocks{device="/dev/nvme0"} 5.00118192e+08
smartctl_device_capacity_blocks{device="/dev/sda"} 2.3437770752e+10
smartctl_device_capacity_blocks{device="/dev/sdb"} 2.3437770752e+10
smartctl_device_capacity_blocks{device="/dev/sdc"} 2.3437770752e+10
smartctl_device_capacity_blocks{device="/dev/sdd"} 2.50069679e+08
# HELP smartctl_device_capacity_bytes Device capacity in bytes
# TYPE smartctl_device_capacity_bytes gauge
smartctl_device_capacity_bytes{device="/dev/nvme0"} 2.56060514304e+11
smartctl_device_capacity_bytes{device="/dev/sda"} 1.2000138625024e+13
smartctl_device_capacity_bytes{device="/dev/sdb"} 1.2000138625024e+13
smartctl_device_capacity_bytes{device="/dev/sdc"} 1.2000138625024e+13
smartctl_device_capacity_bytes{device="/dev/sdd"} 1.28035675648e+11
# HELP smartctl_device_critical_warning This field indicates critical warnings for the state of the controller
# TYPE smartctl_device_critical_warning counter
smartctl_device_critical_warning{device="/dev/nvme0"} 4
smartctl_device_critical_warning{device="/dev/sda"} 0
smartctl_device_critical_warning{device="/dev/sdb"} 0
smartctl_device_critical_warning{device="/dev/sdc"} 0
smartctl_device_critical_warning{device="/dev/sdd"} 0
# HELP smartctl_device_interface_speed Device interface speed, bits per second
# TYPE smartctl_device_interface_speed gauge
smartctl_device_interface_speed{device="/dev/nvme0",speed_type="current"} 0
smartctl_device_interface_speed{device="/dev/nvme0",speed_type="max"} 0
smartctl_device_interface_speed{device="/dev/sda",speed_type="current"} 6e+09
smartctl_device_interface_speed{device="/dev/sda",speed_type="max"} 6e+09
smartctl_device_interface_speed{device="/dev/sdb",speed_type="current"} 6e+09
smartctl_device_interface_speed{device="/dev/sdb",speed_type="max"} 6e+09
smartctl_device_interface_speed{device="/dev/sdc",speed_type="current"} 3e+09
smartctl_device_interface_speed{device="/dev/sdc",speed_type="max"} 6e+09
smartctl_device_interface_speed{device="/dev/sdd",speed_type="current"} 3e+09
smartctl_device_interface_speed{device="/dev/sdd",speed_type="max"} 6e+09
# HELP smartctl_device_media_errors Contains the number of occurrences where the controller detected an unrecovered data integrity error. Errors such as uncorrectable ECC, CRC checksum failure, or LBA tag mismatch are included in this field
# TYPE smartctl_device_media_errors counter
smartctl_device_media_errors{device="/dev/nvme0"} 0
smartctl_device_media_errors{device="/dev/sda"} 0
smartctl_device_media_errors{device="/dev/sdb"} 0
smartctl_device_media_errors{device="/dev/sdc"} 0
smartctl_device_media_errors{device="/dev/sdd"} 0
# HELP smartctl_device_num_err_log_entries Contains the number of Error Information log entries over the life of the controller
# TYPE smartctl_device_num_err_log_entries counter
smartctl_device_num_err_log_entries{device="/dev/nvme0"} 0
smartctl_device_num_err_log_entries{device="/dev/sda"} 0
smartctl_device_num_err_log_entries{device="/dev/sdb"} 0
smartctl_device_num_err_log_entries{device="/dev/sdc"} 0
smartctl_device_num_err_log_entries{device="/dev/sdd"} 0
# HELP smartctl_device_percentage_used Device write percentage used
# TYPE smartctl_device_percentage_used counter
smartctl_device_percentage_used{device="/dev/nvme0"} 161
smartctl_device_percentage_used{device="/dev/sda"} 0
smartctl_device_percentage_used{device="/dev/sdb"} 0
smartctl_device_percentage_used{device="/dev/sdc"} 0
smartctl_device_percentage_used{device="/dev/sdd"} 0
# HELP smartctl_device_power_cycle_count Device power cycle count
# TYPE smartctl_device_power_cycle_count counter
smartctl_device_power_cycle_count{device="/dev/nvme0"} 11
smartctl_device_power_cycle_count{device="/dev/sda"} 11
smartctl_device_power_cycle_count{device="/dev/sdb"} 11
smartctl_device_power_cycle_count{device="/dev/sdc"} 11
smartctl_device_power_cycle_count{device="/dev/sdd"} 44
# HELP smartctl_device_power_on_seconds Device power on seconds
# TYPE smartctl_device_power_on_seconds counter
smartctl_device_power_on_seconds{device="/dev/nvme0"} 5.13324e+07
smartctl_device_power_on_seconds{device="/dev/sda"} 4.80924e+07
smartctl_device_power_on_seconds{device="/dev/sdb"} 4.82652e+07
smartctl_device_power_on_seconds{device="/dev/sdc"} 4.83552e+07
smartctl_device_power_on_seconds{device="/dev/sdd"} 1.346832e+08
# HELP smartctl_device_rotation_rate Device rotation rate
# TYPE smartctl_device_rotation_rate gauge
smartctl_device_rotation_rate{device="/dev/sda"} 7200
smartctl_device_rotation_rate{device="/dev/sdb"} 7200
smartctl_device_rotation_rate{device="/dev/sdc"} 7200
# HELP smartctl_device_smart_status General smart status
# TYPE smartctl_device_smart_status gauge
smartctl_device_smart_status{device="/dev/nvme0"} 0
smartctl_device_smart_status{device="/dev/sda"} 1
smartctl_device_smart_status{device="/dev/sdb"} 1
smartctl_device_smart_status{device="/dev/sdc"} 1
smartctl_device_smart_status{device="/dev/sdd"} 1
# HELP smartctl_device_smartctl_exit_status Exit status of smartctl on device
# TYPE smartctl_device_smartctl_exit_status gauge
smartctl_device_smartctl_exit_status{device="/dev/nvme0"} 8
smartctl_device_smartctl_exit_status{device="/dev/sda"} 0
smartctl_device_smartctl_exit_status{device="/dev/sdb"} 0
smartctl_device_smartctl_exit_status{device="/dev/sdc"} 0
smartctl_device_smartctl_exit_status{device="/dev/sdd"} 0
# HELP smartctl_device_status Device status
# TYPE smartctl_device_status gauge
smartctl_device_status{device="/dev/nvme0"} 0
smartctl_device_status{device="/dev/sda"} 1
smartctl_device_status{device="/dev/sdb"} 1
smartctl_device_status{device="/dev/sdc"} 1
smartctl_device_status{device="/dev/sdd"} 1
# HELP smartctl_device_temperature Device temperature celsius
# TYPE smartctl_device_temperature gauge
smartctl_device_temperature{device="/dev/nvme0",temperature_type="current"} 37
smartctl_device_temperature{device="/dev/sda",temperature_type="current"} 40
smartctl_device_temperature{device="/dev/sdb",temperature_type="current"} 40
smartctl_device_temperature{device="/dev/sdc",temperature_type="current"} 39
smartctl_device_temperature{device="/dev/sdd",temperature_type="current"} 39
# HELP smartctl_version smartctl version
# TYPE smartctl_version gauge
smartctl_version{build_info="(local build)",json_format_version="1.0",smartctl_version="7.3",svn_revision="5338"} 1
```
