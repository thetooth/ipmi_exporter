package collector

import (
	"flag"
	"os/exec"
	"strconv"
	"strings"

	// Prometheus Go toolset
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	ipmiTarget string
	ipmiUser   string
	ipmiPasswd string
)

func init() {
	flag.StringVar(&ipmiTarget, "ipmi-target", "", "IPMI target address")
	flag.StringVar(&ipmiUser, "ipmi-user", "admin", "IPMI username")
	flag.StringVar(&ipmiPasswd, "ipmi-passwd", "", "IPMI password")
}

// SDRCollector declares the data type within the prometheus metrics package.
type SDRCollector struct {
	temperature *prometheus.GaugeVec
	voltage     *prometheus.GaugeVec
	fanSpeed    *prometheus.GaugeVec
}

// NewSDRExporter returns a newly allocated exporter SDRCollector.
// It exposes sensors as reported by ipmitool
func NewSDRExporter() (*SDRCollector, error) {
	return &SDRCollector{
		temperature: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ipmi_temperature",
			Help: "Temperature in degrees C.",
		}, []string{"status", "sensor"}),
		voltage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ipmi_voltage",
			Help: "Device voltage.",
		}, []string{"status", "bus"}),
		fanSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ipmi_fan_speed",
			Help: "Fan speed in RPM.",
		}, []string{"status", "fan"}),
	}, nil
}

// Describe describes all the metrics.
func (e *SDRCollector) Describe(ch chan<- *prometheus.Desc) {
	e.temperature.Describe(ch)
	e.voltage.Describe(ch)
	e.fanSpeed.Describe(ch)
}

// Collect fetches the stats.
func (e *SDRCollector) Collect(ch chan<- prometheus.Metric) {
	e.ipmitool()
	e.temperature.Collect(ch)
	e.voltage.Collect(ch)
	e.fanSpeed.Collect(ch)
}

// ipmitool -I lanplus -H 1.1.1.1 -U admin -P **** sdr
func (e *SDRCollector) ipmitool() {
	out, err := exec.Command("ipmitool", "-I", "lanplus", "-H", ipmiTarget, "-U", ipmiUser, "-P", ipmiPasswd, "sdr").Output()
	if err != nil {
		log.Errorf("error on executing ipmitool: %v", err)
		return
	}
	outlines := strings.Split(string(out), "\n")
	for _, line := range outlines {
		parsedLine := strings.FieldsFunc(line, func(r rune) bool {
			return r == '|'
		})

		// Skip lastline/bad output
		if len(parsedLine) != 3 {
			continue
		}
		name := strings.TrimSpace(parsedLine[0])
		raw := strings.TrimSpace(parsedLine[1])
		status := strings.TrimSpace(parsedLine[2])

		// Skip sensors with no reading
		if status == "ns" {
			continue
		}

		if strings.Contains(raw, "Volts") {
			dec := strings.Fields(raw)[0]
			volt, err := strconv.ParseFloat(dec, 64)
			if err != nil {
				log.Errorln(err)
				return
			}
			e.voltage.With(prometheus.Labels{"status": status, "bus": name}).Set(volt)
		} else if strings.Contains(raw, "degrees") {
			dec := strings.Fields(raw)[0]
			celsius, err := strconv.ParseFloat(dec, 64)
			if err != nil {
				log.Errorln(err)
				return
			}
			e.temperature.With(prometheus.Labels{"status": status, "sensor": name}).Set(celsius)
		} else if strings.Contains(raw, "RPM") {
			dec := strings.Fields(raw)[0]
			rpm, err := strconv.ParseFloat(dec, 64)
			if err != nil {
				log.Errorln(err)
				return
			}
			e.fanSpeed.With(prometheus.Labels{"status": status, "fan": name}).Set(rpm)
		} else {
			log.Debug("Skipping due to unknown type")
			continue
		}
	}
}
