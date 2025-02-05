package main

import (
	"flag"
	"net/http"

	"github.com/jankoudela/amd-smi-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	listenAddress = flag.String("web.listen-address", ":9360", "Address to listen on for AMD GPU metrics.")
	runtimeListenAddress = flag.String("web.runtime-listen-address", ":9361", "Address to listen on for Go runtime metrics.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	debug        = flag.Bool("debug", false, "Enable debug logging")
)

func main() {
	flag.Parse()

	// Configure logging
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Set up panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic in main: %v", r)
		}
	}()

	// Create separate registries for AMD GPU and runtime metrics
	amdRegistry := prometheus.NewRegistry()
	runtimeRegistry := prometheus.NewRegistry()

	// Register AMD GPU metrics
	collector := collector.NewAMDSMICollector()
	amdRegistry.MustRegister(collector)

	// Register runtime metrics
	runtimeRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	runtimeRegistry.MustRegister(prometheus.NewGoCollector())

	// AMD GPU metrics server
	amdMux := http.NewServeMux()
	amdMux.Handle(*metricsPath, promhttp.HandlerFor(amdRegistry, promhttp.HandlerOpts{}))
	amdMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>AMD SMI Exporter</title></head>
			<body>
			<h1>AMD SMI Exporter</h1>
			<p><a href="` + *metricsPath + `">AMD GPU Metrics</a></p>
			</body>
			</html>`))
	})

	// Runtime metrics server
	runtimeMux := http.NewServeMux()
	runtimeMux.Handle(*metricsPath, promhttp.HandlerFor(runtimeRegistry, promhttp.HandlerOpts{}))
	runtimeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>AMD SMI Exporter Runtime</title></head>
			<body>
			<h1>AMD SMI Exporter Runtime</h1>
			<p><a href="` + *metricsPath + `">Go Runtime Metrics</a></p>
			</body>
			</html>`))
	})

	// Start servers
	log.Infof("Starting AMD GPU metrics server on %s", *listenAddress)
	log.Infof("Starting runtime metrics server on %s", *runtimeListenAddress)

	// Start runtime metrics server in a goroutine
	go func() {
		if err := http.ListenAndServe(*runtimeListenAddress, runtimeMux); err != nil {
			log.Fatal(err)
		}
	}()

	// Start AMD GPU metrics server
	if err := http.ListenAndServe(*listenAddress, amdMux); err != nil {
		log.Fatal(err)
	}
}
