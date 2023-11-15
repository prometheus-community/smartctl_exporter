This directory contains JSON testing data for parsing validation.

New data can be collected with the `collect_fake_json.sh` script.

Sensitive information should been redacted using the `redact_fake_json.py`
script.

TODO: what is a good naming scheme for files in this directory? For first-pass,
it has been either of `model_name` or `scsi_model_name`, followed by an
identifier. Why multiple drives of the same model? Testing where one drive has
fields not present on others, e.g. error counts.
