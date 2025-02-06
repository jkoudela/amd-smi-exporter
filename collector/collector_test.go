package collector

import (
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/sirupsen/logrus"
)

func TestAMDSMICollectorWithSampleData(t *testing.T) {
	// Save current osExecute and restore it after the test
	originalOsExecute := osExecute
	defer func() { osExecute = originalOsExecute }()





	// Mock the osExecute function to return sample file data
	osExecute = func(command string) ([]byte, error) {
		// Read the sample file
		sampleData, err := os.ReadFile("../amd-smi-output.sample")
		if err != nil {
			return nil, err
		}
		
		// Print the JSON for debugging
		t.Logf("Sample JSON: %s", string(sampleData))
		
		return sampleData, nil
	}

	// Set log level to debug for testing
	logrus.SetLevel(logrus.DebugLevel)

	// Create a new collector
	collector := NewAMDSMICollector()

	// Create a new registry and register the collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Trigger collection
	if err := testutil.CollectAndCompare(collector, strings.NewReader(`
# HELP amd_gpu_clock_mhz GPU clock metrics in MHz
# TYPE amd_gpu_clock_mhz gauge
amd_gpu_clock_mhz{gpu="0",type="gfx"} 60
# HELP amd_gpu_ecc_errors_total GPU ECC error counts
# TYPE amd_gpu_ecc_errors_total gauge
amd_gpu_ecc_errors_total{gpu="0",type="correctable"} 0
amd_gpu_ecc_errors_total{gpu="0",type="uncorrectable"} 0
# HELP amd_gpu_fan GPU fan metrics
# TYPE amd_gpu_fan gauge
amd_gpu_fan{gpu="0",type="rpm"} 981
amd_gpu_fan{gpu="0",type="speed"} 51
amd_gpu_fan{gpu="0",type="usage"} 20
# HELP amd_gpu_memory_bytes GPU memory usage in bytes
# TYPE amd_gpu_memory_bytes gauge
amd_gpu_memory_bytes{gpu="0",type="free_gtt"} 134738870272
amd_gpu_memory_bytes{gpu="0",type="free_visible_vram"} 46458208256
amd_gpu_memory_bytes{gpu="0",type="free_vram"} 46458208256
amd_gpu_memory_bytes{gpu="0",type="total_gtt"} 135053443072
amd_gpu_memory_bytes{gpu="0",type="total_visible_vram"} 48301604864
amd_gpu_memory_bytes{gpu="0",type="total_vram"} 48301604864
amd_gpu_memory_bytes{gpu="0",type="used_gtt"} 314572800
amd_gpu_memory_bytes{gpu="0",type="used_visible_vram"} 1843396608
amd_gpu_memory_bytes{gpu="0",type="used_vram"} 1843396608
# HELP amd_gpu_power_watts GPU power consumption in watts
# TYPE amd_gpu_power_watts gauge
amd_gpu_power_watts{gpu="0",type="socket"} 34
# HELP amd_gpu_temperature_celsius GPU temperature in celsius
# TYPE amd_gpu_temperature_celsius gauge
amd_gpu_temperature_celsius{gpu="0",type="edge"} 43
amd_gpu_temperature_celsius{gpu="0",type="hotspot"} 49
amd_gpu_temperature_celsius{gpu="0",type="memory"} 56
# HELP amd_gpu_usage_percent GPU usage metrics in percent
# TYPE amd_gpu_usage_percent gauge
amd_gpu_usage_percent{gpu="0",type="gfx"} 0
amd_gpu_usage_percent{gpu="0",type="mm"} 0
amd_gpu_usage_percent{gpu="0",type="umc"} 0
amd_gpu_usage_percent{gpu="0",type="vcn"} 19
# HELP amd_gpu_voltage_mv GPU voltage in millivolts
# TYPE amd_gpu_voltage_mv gauge
amd_gpu_voltage_mv{gpu="0",type="gfx"} 729
amd_gpu_voltage_mv{gpu="0",type="memory"} 674
amd_gpu_voltage_mv{gpu="0",type="soc"} 704
`)); err != nil {
		t.Errorf("unexpected collecting result: %s", err)
	}
}
