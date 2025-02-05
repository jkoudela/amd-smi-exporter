# AMD SMI Exporter

A Prometheus exporter for AMD GPU metrics using the `amd-smi` tool.

## Prerequisites

- Go 1.21 or later
- AMD GPU with `amd-smi` tool installed

## Installation

```bash
go get github.com/jankoudela/amd-smi-exporter
```

## Building

The project includes a Makefile with several build targets:

### Build for local architecture (macOS)
```bash
make build
```

### Build for Linux (amd64)
```bash
make build-linux
```
This will create a binary named `amd-smi-exporter-linux-amd64` that you can deploy to your Linux machine.

### Clean build artifacts
```bash
make clean
```

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
- `amd_gpu_power_watts`: GPU power consumption in watts
- `amd_gpu_temperature_celsius`: GPU temperature in celsius
- `amd_gpu_clock_mhz`: GPU clock speeds in MHz
- `amd_gpu_memory_bytes`: GPU memory usage in bytes
- `amd_gpu_fan`: GPU fan metrics (speed and RPM)
- `amd_gpu_voltage_mv`: GPU voltage in millivolts
- `amd_gpu_ecc_errors_total`: GPU ECC error counts

Each metric includes appropriate labels to identify the GPU and metric type.

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'amd_gpu'
    static_configs:
      - targets: ['localhost:9360']
```

