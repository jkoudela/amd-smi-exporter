package collector

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// osExecute is a variable that can be replaced in tests
var osExecute = func(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)
	return cmd.Output()
}

type AMDSMICollector struct {
	gpuUsage       *prometheus.GaugeVec
	gpuPower       *prometheus.GaugeVec
	gpuTemperature *prometheus.GaugeVec
	gpuClock       *prometheus.GaugeVec
	gpuMemoryUsage *prometheus.GaugeVec
	gpuFan         *prometheus.GaugeVec
	gpuVoltage     *prometheus.GaugeVec
	gpuEccErrors   *prometheus.GaugeVec
}

func NewAMDSMICollector() *AMDSMICollector {
	return &AMDSMICollector{
		gpuUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_usage_percent",
				Help: "GPU usage metrics in percent",
			},
			[]string{"gpu", "type"},
		),
		gpuPower: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_power_watts",
				Help: "GPU power consumption in watts",
			},
			[]string{"gpu", "type"},
		),
		gpuTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_temperature_celsius",
				Help: "GPU temperature in celsius",
			},
			[]string{"gpu", "type"},
		),
		gpuClock: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_clock_mhz",
				Help: "GPU clock metrics in MHz",
			},
			[]string{"gpu", "type"},
		),
		gpuMemoryUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_memory_bytes",
				Help: "GPU memory usage in bytes",
			},
			[]string{"gpu", "type"},
		),
		gpuFan: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_fan",
				Help: "GPU fan metrics",
			},
			[]string{"gpu", "type"},
		),
		gpuVoltage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_voltage_mv",
				Help: "GPU voltage in millivolts",
			},
			[]string{"gpu", "type"},
		),
		gpuEccErrors: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "amd_gpu_ecc_errors_total",
				Help: "GPU ECC error counts",
			},
			[]string{"gpu", "type"},
		),
	}
}

func (c *AMDSMICollector) Describe(ch chan<- *prometheus.Desc) {
	c.gpuUsage.Describe(ch)
	c.gpuPower.Describe(ch)
	c.gpuTemperature.Describe(ch)
	c.gpuClock.Describe(ch)
	c.gpuMemoryUsage.Describe(ch)
	c.gpuFan.Describe(ch)
	c.gpuVoltage.Describe(ch)
	c.gpuEccErrors.Describe(ch)
}

func (c *AMDSMICollector) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic in Collect: %v", r)
		}
	}()

	// Execute amd-smi with metrics in JSON format
	command := "COLUMNS=1000 amd-smi metric --json"
	log.Debugf("Executing command: %s", command)

	output, err := osExecute(command)
	if err != nil {
		log.Errorf("Error executing amd-smi: %v, output: %s", err, string(output))
		return
	}

	log.Debugf("Raw JSON output: %s", string(output))

	var data interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		log.Errorf("Error parsing JSON output: %v, raw output: %s", err, string(output))
		return
	}

	log.Debugf("Parsed data type: %T", data)

	c.collectMetrics(data)

	// Collect all metrics
	c.gpuUsage.Collect(ch)
	c.gpuPower.Collect(ch)
	c.gpuTemperature.Collect(ch)
	c.gpuClock.Collect(ch)
	c.gpuMemoryUsage.Collect(ch)
	c.gpuFan.Collect(ch)
	c.gpuVoltage.Collect(ch)
	c.gpuEccErrors.Collect(ch)
}

func (c *AMDSMICollector) collectMetrics(rawData interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic in collectMetrics: %v", r)
		}
	}()

	// Reset all metrics before collecting new values
	c.gpuUsage.Reset()
	c.gpuPower.Reset()
	c.gpuTemperature.Reset()
	c.gpuClock.Reset()
	c.gpuMemoryUsage.Reset()
	c.gpuFan.Reset()
	c.gpuVoltage.Reset()
	c.gpuEccErrors.Reset()

	// Handle array of GPU objects
	gpus, ok := rawData.([]interface{})
	if !ok {
		log.Errorf("Expected array of GPU objects, got %T", rawData)
		return
	}

	// Process each GPU
	for _, gpuData := range gpus {
		gpu, ok := gpuData.(map[string]interface{})
		if !ok {
			log.Errorf("Expected GPU object to be a map, got %T", gpuData)
			continue
		}

		// Get GPU ID
		gpuID := "0"
		if id, ok := gpu["gpu"].(float64); ok {
			gpuID = strconv.FormatFloat(id, 'f', 0, 64)
		}

		c.collectGPUMetrics(gpuID, gpu)
	}
}

func (c *AMDSMICollector) collectGPUMetrics(gpuID string, gpu map[string]interface{}) {
	log.Debugf("Processing GPU %s with data: %+v", gpuID, gpu)

	// Usage metrics
	if usage, ok := gpu["usage"].(map[string]interface{}); ok {
		log.Debugf("Processing usage metrics for GPU %s: %+v", gpuID, usage)
		if gfxActivity, ok := getNestedValue(usage, "gfx_activity", "value"); ok {
			c.gpuUsage.WithLabelValues(gpuID, "gfx").Set(gfxActivity)
		}
		if umcActivity, ok := getNestedValue(usage, "umc_activity", "value"); ok {
			c.gpuUsage.WithLabelValues(gpuID, "umc").Set(umcActivity)
		}
		if mmActivity, ok := getNestedValue(usage, "mm_activity", "value"); ok {
			c.gpuUsage.WithLabelValues(gpuID, "mm").Set(mmActivity)
		}
		// VCN activity (video encoder/decoder)
		if vcnActivity, ok := usage["vcn_activity"].([]interface{}); ok && len(vcnActivity) > 0 {
			if vcn0, ok := vcnActivity[0].(map[string]interface{}); ok {
				if value, ok := parseFloat(vcn0["value"]); ok {
					c.gpuUsage.WithLabelValues(gpuID, "vcn").Set(value)
				}
			}
		}
	}

	// Power metrics
	if power, ok := gpu["power"].(map[string]interface{}); ok {
		if socketPower, ok := getNestedValue(power, "socket_power", "value"); ok {
			c.gpuPower.WithLabelValues(gpuID, "socket").Set(socketPower)
		}
		if gfxVoltage, ok := getNestedValue(power, "gfx_voltage", "value"); ok {
			c.gpuVoltage.WithLabelValues(gpuID, "gfx").Set(gfxVoltage)
		}
		if socVoltage, ok := getNestedValue(power, "soc_voltage", "value"); ok {
			c.gpuVoltage.WithLabelValues(gpuID, "soc").Set(socVoltage)
		}
		if memVoltage, ok := getNestedValue(power, "mem_voltage", "value"); ok {
			c.gpuVoltage.WithLabelValues(gpuID, "memory").Set(memVoltage)
		}
	}

	// Clock metrics
	if clock, ok := gpu["clock"].(map[string]interface{}); ok {
		// GFX clock
		if gfx0, ok := clock["gfx_0"].(map[string]interface{}); ok {
			if clk, ok := getNestedValue(gfx0, "clk", "value"); ok {
				c.gpuClock.WithLabelValues(gpuID, "gfx").Set(clk)
			}
		}
	}

	// Temperature metrics
	if temp, ok := gpu["temperature"].(map[string]interface{}); ok {
		if edge, ok := getNestedValue(temp, "edge", "value"); ok {
			c.gpuTemperature.WithLabelValues(gpuID, "edge").Set(edge)
		}
		if hotspot, ok := getNestedValue(temp, "hotspot", "value"); ok {
			c.gpuTemperature.WithLabelValues(gpuID, "hotspot").Set(hotspot)
		}
		if mem, ok := getNestedValue(temp, "mem", "value"); ok {
			c.gpuTemperature.WithLabelValues(gpuID, "memory").Set(mem)
		}
	}

	// Memory usage metrics
	if memUsage, ok := gpu["mem_usage"].(map[string]interface{}); ok {
		// VRAM metrics
		if value, ok := getNestedValue(memUsage, "total_vram", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "total_vram").Set(value * 1024 * 1024)
		}
		if value, ok := getNestedValue(memUsage, "used_vram", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "used_vram").Set(value * 1024 * 1024)
		}
		if value, ok := getNestedValue(memUsage, "free_vram", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "free_vram").Set(value * 1024 * 1024)
		}

		// Visible VRAM metrics
		if value, ok := getNestedValue(memUsage, "total_visible_vram", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "total_visible_vram").Set(value * 1024 * 1024)
		}
		if value, ok := getNestedValue(memUsage, "used_visible_vram", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "used_visible_vram").Set(value * 1024 * 1024)
		}
		if value, ok := getNestedValue(memUsage, "free_visible_vram", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "free_visible_vram").Set(value * 1024 * 1024)
		}

		// GTT metrics
		if value, ok := getNestedValue(memUsage, "total_gtt", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "total_gtt").Set(value * 1024 * 1024)
		}
		if value, ok := getNestedValue(memUsage, "used_gtt", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "used_gtt").Set(value * 1024 * 1024)
		}
		if value, ok := getNestedValue(memUsage, "free_gtt", "value"); ok {
			c.gpuMemoryUsage.WithLabelValues(gpuID, "free_gtt").Set(value * 1024 * 1024)
		}
	}

	// Fan metrics
	if fan, ok := gpu["fan"].(map[string]interface{}); ok {
		if speed, ok := parseFloat(fan["speed"]); ok {
			c.gpuFan.WithLabelValues(gpuID, "speed").Set(speed)
		}
		if rpm, ok := parseFloat(fan["rpm"]); ok {
			c.gpuFan.WithLabelValues(gpuID, "rpm").Set(rpm)
		}
		if usage, ok := getNestedValue(fan, "usage", "value"); ok {
			c.gpuFan.WithLabelValues(gpuID, "usage").Set(usage)
		}
	}

	// ECC metrics
	if ecc, ok := gpu["ecc"].(map[string]interface{}); ok {
		if correctable, ok := parseFloat(ecc["total_correctable_count"]); ok {
			c.gpuEccErrors.WithLabelValues(gpuID, "correctable").Set(correctable)
		}
		if uncorrectable, ok := parseFloat(ecc["total_uncorrectable_count"]); ok {
			c.gpuEccErrors.WithLabelValues(gpuID, "uncorrectable").Set(uncorrectable)
		}
	}
}

// Helper function to get map keys for debugging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func parseFloat(value interface{}) (float64, bool) {
	if value == nil {
		return 0, false
	}

	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		if v == "N/A" || v == "" {
			return 0, false
		}
		v = strings.TrimSpace(v)
		if strings.HasSuffix(v, " W") {
			v = strings.TrimSuffix(v, " W")
		} else if strings.HasSuffix(v, " MHz") {
			v = strings.TrimSuffix(v, " MHz")
		} else if strings.HasSuffix(v, " %") {
			v = strings.TrimSuffix(v, " %")
		} else if strings.HasSuffix(v, " °C") {
			v = strings.TrimSuffix(v, " °C")
		} else if strings.HasSuffix(v, " MB") {
			v = strings.TrimSuffix(v, " MB")
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		log.Debugf("Could not parse value as float: %v", v)
	}
	return 0, false
}

func getNestedValue(m map[string]interface{}, keys ...string) (float64, bool) {
	log.Debugf("Getting nested value for keys: %v", keys)
	current := m
	for _, key := range keys {
		log.Debugf("Looking for key %s in %+v", key, current)
		if value, ok := current[key]; ok {
			log.Debugf("Found value for key %s: %+v", key, value)
			if nextMap, ok := value.(map[string]interface{}); ok {
				log.Debugf("Value is a map, updating current to: %+v", nextMap)
				current = nextMap
			} else if v, ok := value.(float64); ok {
				return v, true
			} else {
				return 0, false
			}
		} else {
			return 0, false
		}
	}

	// If we reach here, we've traversed all keys and the last value is a map
	if value, ok := current["value"]; ok {
		return parseFloat(value)
	}
	return 0, false
}
