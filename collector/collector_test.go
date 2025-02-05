package collector

import (
	"encoding/json"
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



	// Convert sample data to JSON format
	gpuData := []interface{}{
		map[string]interface{}{
			"gpu": float64(0),
			"usage": map[string]interface{}{
				"gfx_activity": map[string]interface{}{
					"value": float64(0),
				},
				"umc_activity": map[string]interface{}{
					"value": float64(2),
				},
				"mm_activity": map[string]interface{}{
					"value": float64(98),
				},
			},
			"power": map[string]interface{}{
				"socket_power": map[string]interface{}{
					"value": float64(41),
				},
				"gfx_voltage": map[string]interface{}{
					"value": float64(719),
				},
			},
			"temperature": map[string]interface{}{
				"edge": map[string]interface{}{
					"value": float64(42),
				},
				"hotspot": map[string]interface{}{
					"value": float64(49),
				},
			},
			"clock": map[string]interface{}{
				"gfx_0": map[string]interface{}{
					"clk": map[string]interface{}{
						"value": float64(252),
					},
				},
			},
			"memory": map[string]interface{}{
				"used": map[string]interface{}{
					"value": float64(1758 * 1024 * 1024), // Convert MB to bytes
				},
				"total": map[string]interface{}{
					"value": float64(46064 * 1024 * 1024), // Convert MB to bytes
				},
			},
			"fan": map[string]interface{}{
				"speed": map[string]interface{}{
					"value": float64(51),
				},
				"rpm": map[string]interface{}{
					"value": float64(977),
				},
			},
		},
	}

	// Mock the osExecute function to return JSON data
	osExecute = func(command string) ([]byte, error) {
		// The real command would return JSON, so we need to marshal our data
		jsonData, err := json.Marshal(gpuData)
		if err != nil {
			return nil, err
		}
		
		// Print the JSON for debugging
		t.Logf("Generated JSON: %s", string(jsonData))
		
		return jsonData, nil
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
# HELP amd_gpu_usage_percent GPU usage metrics in percent
# TYPE amd_gpu_usage_percent gauge
amd_gpu_usage_percent{gpu="0",type="gfx"} 0
amd_gpu_usage_percent{gpu="0",type="umc"} 2
amd_gpu_usage_percent{gpu="0",type="mm"} 98
# HELP amd_gpu_power_watts GPU power consumption in watts
# TYPE amd_gpu_power_watts gauge
amd_gpu_power_watts{gpu="0",type="socket"} 41
# HELP amd_gpu_temperature_celsius GPU temperature in celsius
# TYPE amd_gpu_temperature_celsius gauge
amd_gpu_temperature_celsius{gpu="0",type="edge"} 42
amd_gpu_temperature_celsius{gpu="0",type="hotspot"} 49
# HELP amd_gpu_clock_mhz GPU clock metrics in MHz
# TYPE amd_gpu_clock_mhz gauge
amd_gpu_clock_mhz{gpu="0",type="gfx"} 252
# HELP amd_gpu_memory_bytes GPU memory usage in bytes
# TYPE amd_gpu_memory_bytes gauge
amd_gpu_memory_bytes{gpu="0",type="total"} 4.8301604864e+10
amd_gpu_memory_bytes{gpu="0",type="used"} 1.843396608e+09
# HELP amd_gpu_voltage_mv GPU voltage in millivolts
# TYPE amd_gpu_voltage_mv gauge
amd_gpu_voltage_mv{gpu="0",type="gfx"} 719
`)); err != nil {
		t.Errorf("unexpected collecting result: %s", err)
	}
}
