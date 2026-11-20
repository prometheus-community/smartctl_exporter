JSONNET_FMT := jsonnetfmt -n 2 --max-blank-lines 2 --string-style s --comment-style s
OUTPUT_DIR := output

all: fmt $(OUTPUT_DIR)/smartctl_alerts.yaml $(OUTPUT_DIR)/dashboards lint

fmt:
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		xargs -n 1 -- $(JSONNET_FMT) -i

$(OUTPUT_DIR)/smartctl_alerts.yaml: mixin.libsonnet config.libsonnet $(wildcard alerts/*)
	@mkdir -p $(OUTPUT_DIR)
	jsonnet -S alerts.jsonnet > $@

$(OUTPUT_DIR)/dashboards: mixin.libsonnet config.libsonnet $(wildcard dashboards/*)
	@mkdir -p $(OUTPUT_DIR)
	jsonnet -m $(OUTPUT_DIR) dashboards.jsonnet

lint: $(OUTPUT_DIR)/smartctl_alerts.yaml
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		while read f; do \
			$(JSONNET_FMT) "$$f" | diff -u "$$f" -; \
		done
	@if command -v promtool >/dev/null 2>&1; then \
		promtool check rules $(OUTPUT_DIR)/smartctl_alerts.yaml; \
	else \
		echo "promtool not found, skipping rules check"; \
	fi

.PHONY: jb_install
jb_install:
	jb install

clean:
	rm -rf $(OUTPUT_DIR)
