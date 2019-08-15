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
		json := readData(device)
		info.SetJSON(json)
		smart := NewSMARTctl(json, ch)
		smart.Collect()
	}
	info.Collect()
}

func init() {
	options = loadOptions()
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
