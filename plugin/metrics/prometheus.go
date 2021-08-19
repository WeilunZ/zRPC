package metrics

import (
	"github.com/WeilunZ/zRPC/plugin"
	"net/http"

	"github.com/WeilunZ/zRPC/components/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "zRPC"
)

var (
	Port            = ":9628"
	Endpoint        = "/metrics"
	DefaultRegistry = newDefaultMetricsRegistry()
)

type Prometheus struct {
	opts *plugin.Options
	Port string
	Endpoint string
	registry *prometheus.Registry
}

func Run() error {
	log.Infof("Starting http server to serve metrics at port '%s', endpoint '%s'", Port, Endpoint)

	server := http.NewServeMux()
	server.Handle(Endpoint, promhttp.HandlerFor(DefaultRegistry, promhttp.HandlerOpts{}))

	// start an http server using the mux server
	return http.ListenAndServe(Port, server)
}

func NewCounter(metricName string) prometheus.Counter {
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      metricName,
	})
	if err := DefaultRegistry.Register(c); err != nil {
		log.Warningf("metric register err: %s", err)
	}
	return c
}

func NewCounterVec(metricName string, labels ...string) *prometheus.CounterVec {
	cv := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      metricName,
	}, labels)
	if err := DefaultRegistry.Register(cv); err != nil {
		log.Warningf("metric register err: %s", err)
	}
	return cv
}

func NewGauge(metricName string) prometheus.Gauge {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      metricName,
	})
	if err := DefaultRegistry.Register(g); err != nil {
		log.Warningf("metric register err: %s", err)
	}
	return g
}

func newDefaultMetricsRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	return registry
}
