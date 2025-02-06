#!/bin/bash
# Installs amd-smi-exporter binary and service locally

set -e

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root"
    exit 1
fi

# Stop existing service if running
if systemctl is-active --quiet amd-smi-exporter; then
    echo "Stopping existing service..."
    systemctl stop amd-smi-exporter
fi

# Install binary
echo "Installing amd-smi-exporter binary..."
install -m 755 amd-smi-exporter /usr/local/bin/

# Install service file
echo "Installing systemd service..."
install -m 644 amd-smi-exporter.service /etc/systemd/system/

# Reload systemd and enable service
echo "Enabling and starting service..."
systemctl daemon-reload
systemctl enable amd-smi-exporter
systemctl start amd-smi-exporter

# Show service status
echo "\nService status:"
systemctl status amd-smi-exporter --no-pager

echo "\nAMD SMI Exporter metrics available at:"
echo "  http://localhost:9360/metrics (GPU metrics)"
echo "  http://localhost:9361/metrics (Runtime metrics)"
