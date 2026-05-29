#!/usr/bin/env bash
set -euo pipefail

# Deploy smartctl_exporter with FARM log support and Grafana dashboard.
#
# Usage:
#   ./deploy.sh <target_host> [target_host2] ...
#
# Examples:
#   ./deploy.sh 192.168.5.114
#   ./deploy.sh 192.168.5.114 192.168.5.115 192.168.5.116
#
# Environment variables:
#   SSH_USER        - SSH user (default: root)
#   GRAFANA_USER    - Grafana admin user (default: admin)
#   GRAFANA_PASS    - Grafana admin password (default: admin)
#   GRAFANA_HOST    - Override Grafana host (default: auto-detect from target)
#   PROMETHEUS_HOST - Host running central Prometheus (default: first target)
#
# Notes:
#   - The Prometheus scrape config is REPLACED with all targets from this run.
#     Always pass ALL hosts you want monitored, not just new ones.
#   - Requires smartmontools >= 7.4 on target nodes for FARM log support.
#
# Prerequisites:
#   - SSH root access to target hosts
#   - Go installed locally for building
#   - Prometheus with scrape_config_files including /etc/prometheus/scrape_configs/

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SSH_USER="${SSH_USER:-root}"
GRAFANA_USER="${GRAFANA_USER:-admin}"
GRAFANA_PASS="${GRAFANA_PASS:-admin}"
GRAFANA_HOST="${GRAFANA_HOST:-}"
PROMETHEUS_HOST="${PROMETHEUS_HOST:-}"

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 <target_host> [target_host2] ..."
    echo ""
    echo "NOTE: Pass ALL monitored hosts each run (scrape config is replaced, not appended)."
    exit 1
fi

echo "==> Building smartctl_exporter (static)..."
cd "$SCRIPT_DIR"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o smartctl_exporter .

for TARGET in "$@"; do
    echo ""
    echo "========================================"
    echo "==> Deploying to ${TARGET}"
    echo "========================================"

    echo "==> Ensuring smartmontools is installed..."
    ssh "${SSH_USER}@${TARGET}" "command -v smartctl >/dev/null 2>&1 || { echo 'Installing smartmontools...'; dnf install -y smartmontools >/dev/null 2>&1 || apt-get install -y smartmontools >/dev/null 2>&1; }"

    echo "==> Checking smartctl version..."
    ssh "${SSH_USER}@${TARGET}" "smartctl --version | head -1"

    echo "==> Stopping existing service (if running)..."
    ssh "${SSH_USER}@${TARGET}" "systemctl stop smartctl_exporter 2>/dev/null || true"

    echo "==> Copying binary..."
    scp smartctl_exporter "${SSH_USER}@${TARGET}:/usr/bin/smartctl_exporter"

    echo "==> Installing systemd service..."
    scp systemd/smartctl_exporter.service "${SSH_USER}@${TARGET}:/etc/systemd/system/smartctl_exporter.service"

    # If smartctl is not at the default path, add --smartctl.path to the service file
    SMARTCTL_PATH=$(ssh "${SSH_USER}@${TARGET}" "which smartctl 2>/dev/null || echo /usr/sbin/smartctl")
    if [[ "$SMARTCTL_PATH" != "/usr/sbin/smartctl" ]]; then
        echo "    smartctl found at ${SMARTCTL_PATH}, updating service file..."
        ssh "${SSH_USER}@${TARGET}" "sed -i 's|--smartctl.farm-log|--smartctl.farm-log --smartctl.path=${SMARTCTL_PATH}|' /etc/systemd/system/smartctl_exporter.service"
    fi

    echo "==> Enabling and starting smartctl_exporter..."
    ssh "${SSH_USER}@${TARGET}" "systemctl daemon-reload && systemctl enable --now smartctl_exporter && systemctl restart smartctl_exporter"

    # Check if Grafana is running on this node and import dashboard
    GRAFANA_TARGET="${GRAFANA_HOST:-$TARGET}"
    if ssh "${SSH_USER}@${TARGET}" "systemctl is-active grafana-server >/dev/null 2>&1"; then
        echo "==> Grafana detected on ${TARGET}, importing dashboard..."
        sleep 5
        # Import dashboard with instance regex cleared (we use IPs, not localhost)
        python3 -c "
import sys, json
with open('${SCRIPT_DIR}/smartctl-farm-dashboard.json') as f:
    d = json.load(f)
for var in d.get('dashboard', d).get('templating', {}).get('list', []):
    if var.get('name') == 'instance':
        var['regex'] = ''
d.setdefault('overwrite', True)
if 'dashboard' in d:
    d['dashboard']['id'] = None
print(json.dumps(d))
" | curl -sf -u "${GRAFANA_USER}:${GRAFANA_PASS}" \
            -X POST -H "Content-Type: application/json" \
            "http://${GRAFANA_TARGET}:3000/api/dashboards/db" \
            -d @- \
            && echo "    Dashboard imported." \
            || echo "    WARNING: Dashboard import failed."
    else
        echo "==> Grafana not running on ${TARGET}, skipping dashboard import."
    fi

    echo "==> ${TARGET} done. Exporter at http://${TARGET}:9633/metrics"
done

# Update central Prometheus scrape config with all targets
PROM_HOST="${PROMETHEUS_HOST:-$1}"
echo ""
echo "==> Updating Prometheus scrape config on ${PROM_HOST}..."

# Check if Prometheus is installed on the target
if ! ssh "${SSH_USER}@${PROM_HOST}" "systemctl list-unit-files prometheus.service >/dev/null 2>&1"; then
    echo "    WARNING: Prometheus not found on ${PROM_HOST}, skipping scrape config."
    echo ""
    echo "==> Deployment complete! (Prometheus scrape config must be configured manually)"
    exit 0
fi

# Build targets list from all deployed hosts
TARGETS_YAML=""
TARGETS_INLINE=""
for TARGET in "$@"; do
    if [[ -n "$TARGETS_INLINE" ]]; then
        TARGETS_INLINE="${TARGETS_INLINE}, "
    fi
    TARGETS_INLINE="${TARGETS_INLINE}'${TARGET}:9633'"
    TARGETS_YAML="${TARGETS_YAML}
      - targets: ['${TARGET}:9633']"
done

# Determine Prometheus config style: scrape_config_files or inline prometheus.yml
if ssh "${SSH_USER}@${PROM_HOST}" "grep -q 'scrape_config_files' /etc/prometheus/prometheus.yml 2>/dev/null"; then
    echo "    Using scrape_config_files style..."
    ssh "${SSH_USER}@${PROM_HOST}" "mkdir -p /etc/prometheus/scrape_configs && cat > /etc/prometheus/scrape_configs/smartctl_exporter.yml" <<EOF
scrape_configs:
  - job_name: 'smartctl'
    static_configs:
      - targets: [${TARGETS_INLINE}]
EOF
else
    echo "    Using inline prometheus.yml style..."
    # Remove existing smartctl job block if present, then append new one
    ssh "${SSH_USER}@${PROM_HOST}" "
        # Remove old smartctl block (job_name line + following lines until next job or EOF)
        sed -i '/job_name:.*smartctl/,/^  - job_name\\|^\$/{/^  - job_name.*smartctl/d;/^  - job_name/!d}' /etc/prometheus/prometheus.yml 2>/dev/null || true
        # Append new smartctl job
        cat >> /etc/prometheus/prometheus.yml <<INNER
  - job_name: 'smartctl'
    scrape_interval: 60s
    static_configs:
      - targets: [${TARGETS_INLINE}]
INNER
    "
fi

echo "==> Restarting Prometheus on ${PROM_HOST}..."
ssh "${SSH_USER}@${PROM_HOST}" "systemctl restart prometheus 2>/dev/null || true"

echo ""
echo "==> Deployment complete!"
