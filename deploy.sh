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
PROMETHEUS_HOST="${PROMETHEUS_HOST:-}"

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 <target_host> [target_host2] ..."
    echo ""
    echo "NOTE: Pass ALL monitored hosts each run (scrape config is replaced, not appended)."
    exit 1
fi

echo "==> Building smartctl_exporter (static)..."
cd "$SCRIPT_DIR"
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o smartctl_exporter .

for TARGET in "$@"; do
    echo ""
    echo "========================================"
    echo "==> Deploying to ${TARGET}"
    echo "========================================"

    echo "==> Ensuring smartmontools >= 7.4 is installed..."
    ssh "${SSH_USER}@${TARGET}" bash -s <<'SMARTCTL_CHECK'
        if ! command -v smartctl >/dev/null 2>&1; then
            echo "    Installing smartmontools..."
            dnf install -y smartmontools 2>/dev/null || apt-get install -y smartmontools 2>/dev/null
        fi
        # Check version and update if below 7.4
        SMARTCTL_VER=$(smartctl --version | head -1 | grep -oP '\d+\.\d+' | head -1)
        MAJOR=$(echo "$SMARTCTL_VER" | cut -d. -f1)
        MINOR=$(echo "$SMARTCTL_VER" | cut -d. -f2)
        if [[ "$MAJOR" -lt 7 ]] || { [[ "$MAJOR" -eq 7 ]] && [[ "$MINOR" -lt 4 ]]; }; then
            echo "    smartmontools ${SMARTCTL_VER} is below 7.4, attempting update..."
            if command -v dnf >/dev/null 2>&1; then
                dnf update -y smartmontools
            elif command -v apt-get >/dev/null 2>&1; then
                apt-get update && apt-get install -y --only-upgrade smartmontools
            fi
            SMARTCTL_VER=$(smartctl --version | head -1 | grep -oP '\d+\.\d+' | head -1)
            MAJOR=$(echo "$SMARTCTL_VER" | cut -d. -f1)
            MINOR=$(echo "$SMARTCTL_VER" | cut -d. -f2)
            if [[ "$MAJOR" -lt 7 ]] || { [[ "$MAJOR" -eq 7 ]] && [[ "$MINOR" -lt 4 ]]; }; then
                echo "    Distro package too old (${SMARTCTL_VER}), building smartmontools 7.4 from source..."
                # Install build dependencies
                if command -v dnf >/dev/null 2>&1; then
                    dnf install -y gcc gcc-c++ make automake tar gzip
                elif command -v apt-get >/dev/null 2>&1; then
                    apt-get update && apt-get install -y build-essential automake tar gzip
                fi
                cd /tmp
                curl -sLO https://sourceforge.net/projects/smartmontools/files/smartmontools/7.4/smartmontools-7.4.tar.gz/download
                mv download smartmontools-7.4.tar.gz
                tar xzf smartmontools-7.4.tar.gz
                cd smartmontools-7.4
                ./configure --prefix=/usr --sysconfdir=/etc
                make -j"$(nproc)"
                make install
                cd /tmp && rm -rf smartmontools-7.4 smartmontools-7.4.tar.gz
                SMARTCTL_VER=$(smartctl --version | head -1 | grep -oP '\d+\.\d+' | head -1)
                echo "    Built and installed smartmontools ${SMARTCTL_VER} from source."
            else
                echo "    Updated to smartmontools ${SMARTCTL_VER}."
            fi
        else
            echo "    smartmontools ${SMARTCTL_VER} OK."
        fi
SMARTCTL_CHECK

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

    # Check if the target has ID_VDEV udev data (ZFS vdev labels / enclosure slots)
    if ssh "${SSH_USER}@${TARGET}" "udevadm info --export --query=property /dev/\$(lsblk -dn -o NAME | head -1) 2>/dev/null | grep -q ID_VDEV"; then
        echo "    ID_VDEV detected, ensuring --smartctl.vdev-label is set..."
        ssh "${SSH_USER}@${TARGET}" "grep -q 'vdev-label' /etc/systemd/system/smartctl_exporter.service || sed -i 's|--smartctl.farm-log|--smartctl.farm-log --smartctl.vdev-label|' /etc/systemd/system/smartctl_exporter.service"
    fi

    echo "==> Enabling and starting smartctl_exporter..."
    ssh "${SSH_USER}@${TARGET}" "systemctl daemon-reload && systemctl enable --now smartctl_exporter && systemctl restart smartctl_exporter"

    echo "==> Opening port 9633/tcp..."
    ssh "${SSH_USER}@${TARGET}" "
        if command -v firewall-cmd >/dev/null 2>&1; then
            firewall-cmd --permanent --add-port=9633/tcp && firewall-cmd --reload
        elif command -v ufw >/dev/null 2>&1; then
            ufw allow 9633/tcp
        else
            echo '    No firewall manager found, skipping.'
        fi
    "

    # Install dashboard via file provisioning if Grafana is on this node
    if ssh "${SSH_USER}@${TARGET}" "systemctl is-active grafana-server >/dev/null 2>&1"; then
        echo "==> Grafana detected on ${TARGET}, installing dashboard to /etc/grafana/dashboards/system/..."
        ssh "${SSH_USER}@${TARGET}" "mkdir -p /etc/grafana/dashboards/system"
        scp "${SCRIPT_DIR}/smartctl-farm-dashboard.json" "${SSH_USER}@${TARGET}:/etc/grafana/dashboards/system/smartctl-farm-dashboard.json"
        ssh "${SSH_USER}@${TARGET}" "chown grafana:grafana /etc/grafana/dashboards/system/smartctl-farm-dashboard.json && chmod 472 /etc/grafana/dashboards/system/smartctl-farm-dashboard.json"
        echo "    Dashboard installed."
    else
        echo "==> Grafana not running on ${TARGET}, skipping dashboard install."
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
    # Use Python for YAML-aware insertion: removes any existing smartctl job,
    # then inserts the new one before the 'alerting:' block (or appends to
    # scrape_configs if no alerting block exists). This avoids the bug where
    # blindly appending to EOF places the job inside the alerting block.
    ssh "${SSH_USER}@${PROM_HOST}" "python3 - /etc/prometheus/prometheus.yml ${TARGETS_INLINE}" <<'PYEOF'
import sys, re

cfg = sys.argv[1]
targets = sys.argv[2:]

with open(cfg) as f:
    content = f.read()

# Remove any existing smartctl job block (idempotent).
content = re.sub(
    r"(?m)^  - job_name: ['\"]smartctl['\"].*\n(?:    [^\n]*\n)*",
    '',
    content
)

# Build the new job block.
target_list = ', '.join(targets)
job = (
    "  - job_name: 'smartctl'\n"
    "    scrape_interval: 60s\n"
    "    static_configs:\n"
    f"      - targets: [{target_list}]\n"
)

# Insert before 'alerting:' if present, otherwise append to end.
if re.search(r'^alerting:', content, re.MULTILINE):
    content = re.sub(
        r'^(alerting:)',
        lambda m: job + m.group(1),
        content,
        count=1,
        flags=re.MULTILINE
    )
else:
    content = content.rstrip('\n') + '\n' + job

with open(cfg, 'w') as f:
    f.write(content)
PYEOF
fi

echo "==> Restarting Prometheus on ${PROM_HOST}..."
ssh "${SSH_USER}@${PROM_HOST}" "systemctl restart prometheus 2>/dev/null || true"

echo ""
echo "==> Deployment complete!"
