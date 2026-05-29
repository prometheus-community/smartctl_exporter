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

echo "==> Building smartctl_exporter..."
cd "$SCRIPT_DIR"
GOOS=linux GOARCH=amd64 go build -o smartctl_exporter .

for TARGET in "$@"; do
    echo ""
    echo "========================================"
    echo "==> Deploying to ${TARGET}"
    echo "========================================"

    echo "==> Ensuring smartmontools is installed..."
    ssh "${SSH_USER}@${TARGET}" "command -v smartctl >/dev/null 2>&1 || { echo 'Installing smartmontools...'; dnf install -y smartmontools >/dev/null 2>&1 || apt-get install -y smartmontools >/dev/null 2>&1; }"

    echo "==> Checking smartctl version..."
    ssh "${SSH_USER}@${TARGET}" "smartctl --version | head -1"

    echo "==> Copying binary..."
    scp smartctl_exporter "${SSH_USER}@${TARGET}:/usr/bin/smartctl_exporter"

    echo "==> Installing systemd service..."
    scp systemd/smartctl_exporter.service "${SSH_USER}@${TARGET}:/etc/systemd/system/smartctl_exporter.service"

    echo "==> Enabling and starting smartctl_exporter..."
    ssh "${SSH_USER}@${TARGET}" "systemctl daemon-reload && systemctl enable --now smartctl_exporter && systemctl restart smartctl_exporter"

    # Check if Grafana is running on this node and import dashboard
    GRAFANA_TARGET="${GRAFANA_HOST:-$TARGET}"
    if ssh "${SSH_USER}@${TARGET}" "systemctl is-active grafana-server >/dev/null 2>&1"; then
        echo "==> Grafana detected on ${TARGET}, importing dashboard..."
        sleep 5
        curl -sf -u "${GRAFANA_USER}:${GRAFANA_PASS}" \
            -X POST -H "Content-Type: application/json" \
            "http://${GRAFANA_TARGET}:3000/api/dashboards/db" \
            -d @"${SCRIPT_DIR}/smartctl-farm-dashboard.json" \
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
TARGETS=""
for TARGET in "$@"; do
    if [[ -n "$TARGETS" ]]; then
        TARGETS="${TARGETS}, "
    fi
    TARGETS="${TARGETS}'${TARGET}:9633'"
done

ssh "${SSH_USER}@${PROM_HOST}" "mkdir -p /etc/prometheus/scrape_configs && cat > /etc/prometheus/scrape_configs/smartctl_exporter.yml" <<EOF
scrape_configs:
  - job_name: 'smartctl'
    static_configs:
      - targets: [${TARGETS}]
EOF

echo "==> Restarting Prometheus on ${PROM_HOST}..."
ssh "${SSH_USER}@${PROM_HOST}" "systemctl restart prometheus 2>/dev/null || true"

echo ""
echo "==> Deployment complete!"
