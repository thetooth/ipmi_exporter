package main

import (
	"flag"
	"net/http"

	"github.com/thetooth/ipmi_exporter/collector"

	// Prometheus Go toolset
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"

	"github.com/kouhin/envflag"
)

var (
	listenAddress string
)

func init() {
	flag.StringVar(&listenAddress, "listen-address", "0.0.0.0:9100", "Address on which to expose metrics and web interface.")
}

// program starter
func main() {
	if err := envflag.Parse(); err != nil {
		panic(err)
	}

	log.Infoln("Starting ipmi_exporter", version.Info())

	// Environmental sensor metrics
	sdr, _ := collector.NewSDRExporter()
	prometheus.MustRegister(sdr)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Infoln("Listening on", listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
