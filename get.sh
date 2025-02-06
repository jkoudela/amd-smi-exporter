#!/bin/bash
# Quick installer for amd-smi-exporter
set -e

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root"
    exit 1
fi

echo "Installing AMD SMI Exporter..."

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Get latest release version
LATEST_VERSION=$(curl -s https://api.github.com/repos/jkoudela/amd-smi-exporter/releases/latest | grep -Po '"tag_name": "\K.*?(?=")')
if [ -z "$LATEST_VERSION" ]; then
    echo "Error: Could not determine latest version"
    exit 1
fi

echo "Downloading version ${LATEST_VERSION}..."

# Download files
curl -L -O "https://github.com/jkoudela/amd-smi-exporter/releases/download/${LATEST_VERSION}/amd-smi-exporter"
curl -L -O "https://github.com/jkoudela/amd-smi-exporter/releases/download/${LATEST_VERSION}/amd-smi-exporter.service"

# Stop existing service if running
if systemctl is-active --quiet amd-smi-exporter; then
    echo "Stopping existing service..."
    systemctl stop amd-smi-exporter
fi

# Install binary
echo "Installing binary..."
install -m 755 amd-smi-exporter /usr/local/bin/

# Install service file
echo "Installing systemd service..."
install -m 644 amd-smi-exporter.service /etc/systemd/system/

# Reload systemd and enable service
echo "Enabling and starting service..."
systemctl daemon-reload
systemctl enable amd-smi-exporter
systemctl start amd-smi-exporter

# Clean up
cd - > /dev/null
rm -rf "$TMP_DIR"

# Show service status
echo -e "\nService status:"
systemctl status amd-smi-exporter --no-pager

echo -e "\nAMD SMI Exporter metrics available at:"
echo "  http://localhost:9360/metrics (GPU metrics)"
echo "  http://localhost:9361/metrics (Runtime metrics)"
