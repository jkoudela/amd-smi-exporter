# AMD SMI Exporter

A Prometheus exporter for AMD GPU metrics using the `amd-smi` tool. The goal was to get some AMD GPU visualization to Grafana.

## Grafana Dashboards

![Dashboard](image.png)

The project includes two Grafana dashboards:

### Main Dashboard (`grafana/amd_gpu_dashboard.json`)

A monitoring dashboard for AMD GPUs with the following features:

- Real-time GPU metrics
- Dark theme optimized visualization
- Multiple panels tracking various GPU metrics
- Tagged with 'gpu' and 'amd' for easy discovery
- Prometheus data source integration

### Metrics Dashboard (`grafana/amd_gpu_metrics.json`)

A detailed metrics-focused dashboard that provides:
- Detailed fan metrics (RPM, Speed, Usage)
- Power consumption metrics in a dedicated row
- Advanced performance metrics

### Dashboard Setup

1. Copy `.env.example` to `.env` and configure your Grafana settings:
   ```bash
   GRAFANA_URL=http://your-grafana-url:3000
   GRAFANA_API_KEY=your-api-key
   ```

2. Use the provided Makefile targets to update the dashboards:
   ```bash
   make update-dashboard  # Updates the main dashboard
   make update-metrics    # Updates the metrics dashboard
   ```
   
Both dashboards will automatically start displaying metrics from your AMD GPUs once they're imported into Grafana.

## Prerequisites

- Go 1.21 or later
- AMD GPU with `amd-smi` tool installed

## Installation

### Option 1: Quick Install (Recommended)

Install the latest release with a single command:

```bash
curl -sSL https://raw.githubusercontent.com/jkoudela/amd-smi-exporter/main/get.sh | sudo bash
```

This will automatically:
- Download the latest release
- Install the binary and service
- Start and enable the service
- Show the service status and available endpoints

### Option 2: Manual Install

If you prefer to inspect the files before installation:

```bash
# Download installation files
curl -L -O "https://github.com/jkoudela/amd-smi-exporter/releases/latest/download/amd-smi-exporter"
curl -L -O "https://github.com/jkoudela/amd-smi-exporter/releases/latest/download/amd-smi-exporter.service"
curl -L -O "https://github.com/jkoudela/amd-smi-exporter/releases/latest/download/install_local.sh"

# Make scripts executable
chmod +x amd-smi-exporter install_local.sh

# Install binary and service (requires sudo)
sudo ./install_local.sh
```

Verify the service is running:
```bash
systemctl status amd-smi-exporter
```

### Option 3: Build from Source

Clone the repository and build the binary:

```bash
git clone https://github.com/jkoudela/amd-smi-exporter.git
cd amd-smi-exporter
make build
```

Alternatively, you can use the provided install script to deploy to a remote Linux host with IP and username using ssh:

```bash
./install_remote.sh <target-host> root
```

## Building

The project includes a Makefile with the following targets:

### Build for Linux (amd64)
```bash
make build
```
Creates a Linux binary named `amd-smi-exporter` that you can deploy to your Linux machine.

### Clean build artifacts
```bash
make clean
```
Removes build artifacts and temporary files.

### Update Grafana Dashboards
```bash
make update-dashboard  # Updates the main dashboard
make update-metrics    # Updates the metrics dashboard
```
Updates the respective Grafana dashboards using the configuration from your `.env` file.

## Usage

```bash
./amd-smi-exporter [flags]
```

### Flags

- `--web.listen-address`: Address to listen on for web interface and telemetry (default ":9360")
- `--web.telemetry-path`: Path under which to expose metrics (default "/metrics")

## Metrics

The exporter provides the following metrics:

- `amd_gpu_usage_percent`: GPU usage metrics in percent
  - Graphics engine usage (gfx)
  - Memory controller usage (umc)
  - Multimedia engine usage (mm)
  - Video codec engine usage (vcn)
- `amd_gpu_power_watts`: GPU power consumption in watts
  - Socket power consumption
- `amd_gpu_temperature_celsius`: GPU temperature in celsius
  - Edge temperature
  - Hotspot temperature
  - Memory temperature
- `amd_gpu_clock_mhz`: GPU clock speeds in MHz
  - Graphics engine clock (gfx)
- `amd_gpu_memory_bytes`: GPU memory usage in bytes
  - Reports VRAM (total, free, used)
  - Reports visible VRAM (total, free, used)
  - Reports GTT memory (total, free, used)
- `amd_gpu_fan`: GPU fan metrics
  - Speed percentage
  - RPM
  - Usage percentage
- `amd_gpu_voltage_mv`: GPU voltage in millivolts
  - Graphics engine voltage (gfx)
  - Memory voltage
  - SOC voltage
- `amd_gpu_ecc_errors_total`: GPU ECC error counts
  - Correctable errors
  - Uncorrectable errors

Each metric includes appropriate labels to identify the GPU and metric type.

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'amd_gpu'
    static_configs:
      - targets: ['localhost:9360']
```

