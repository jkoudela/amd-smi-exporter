#!/bin/bash
# Installs the binary and service file to a remote host Linux host

# Check if target host is provided
if [ $# -lt 1 ]; then
    echo "Usage: $0 <remote_host> [remote_user]"
    echo "Example: $0 server.example.com ubuntu"
    exit 1
fi

REMOTE_HOST=$1
REMOTE_USER=${2:-$(whoami)}  # Use current user if not specified

# Build the binary
echo "Building amd-smi-exporter..."
GOOS=linux GOARCH=amd64 go build -o amd-smi-exporter

# Copy files to remote host
echo "Copying files to ${REMOTE_USER}@${REMOTE_HOST}..."
scp amd-smi-exporter "${REMOTE_USER}@${REMOTE_HOST}:/tmp/"
scp amd-smi-exporter.service "${REMOTE_USER}@${REMOTE_HOST}:/tmp/"

# Install on remote host
echo "Installing on remote host..."
ssh "${REMOTE_USER}@${REMOTE_HOST}" '
    # Stop existing service if running
    sudo systemctl stop amd-smi-exporter.service || true

    # Move binary and set permissions
    sudo mv /tmp/amd-smi-exporter /usr/local/bin/
    sudo chmod 755 /usr/local/bin/amd-smi-exporter

    # Install systemd service
    sudo mv /tmp/amd-smi-exporter.service /etc/systemd/system/
    sudo chmod 644 /etc/systemd/system/amd-smi-exporter.service

    # Reload systemd and start service
    sudo systemctl daemon-reload
    sudo systemctl enable amd-smi-exporter.service
    sudo systemctl start amd-smi-exporter.service

    # Show status
    echo "Service status:"
    sudo systemctl status amd-smi-exporter.service
'

echo "Installation complete!"
echo "AMD SMI Exporter metrics available at:"
echo "  http://${REMOTE_HOST}:9360/metrics (GPU metrics)"
echo "  http://${REMOTE_HOST}:9361/metrics (Runtime metrics)"
