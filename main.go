package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	options Options
	logger  Logger
)

// SMARTctlManagerCollector implements the Collector interface.
type SMARTctlManagerCollector struct {
}

// Describe sends the super-set of all possible descriptors of metrics
func (i SMARTctlManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(i, ch)
}

// Collect is called by the Prometheus registry when collecting metrics.
func (i SMARTctlManagerCollector) Collect(ch chan<- prometheus.Metric) {
	info := NewSMARTctlInfo(ch)
	for _, device := range options.SMARTctl.Devices {
		if json, err := readData(device); err == nil {
			info.SetJSON(json)
			smart := NewSMARTctl(json, ch)
			smart.Collect()
		} else {
			logger.Error(err.Error())
		}
	}
	info.Collect()
}

func init() {
	options = loadOptions()

	if len(options.SMARTctl.Devices) == 0 {
		logger.Debug("No devices specified, trying to load them automatically")
		json := scanSMARTctlDevices()
		devices := json.Get("devices").Array()
		for _, d := range devices {
			device := d.Get("name").String()
			logger.Debug("Found device: %s", device)
			options.SMARTctl.Devices = append(options.SMARTctl.Devices, device)
		}
	}
}

func main() {
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	prometheus.WrapRegistererWithPrefix("", reg).MustRegister(SMARTctlManagerCollector{})

	logger.Info("Starting on %s%s", options.SMARTctl.BindTo, options.SMARTctl.URLPath)
	http.Handle(options.SMARTctl.URLPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(options.SMARTctl.BindTo, nil))
}
