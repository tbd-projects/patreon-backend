package prometheus_monitoring

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	Hits *prometheus.CounterVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	metrics := &PrometheusMetrics{
		Hits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hits",
			Help: "Count requests on service",
		}, []string{"status", "path", "method"}),
	}

	return metrics
}

func (p *PrometheusMetrics) GetHits() *prometheus.CounterVec {
	return p.Hits
}

func (p *PrometheusMetrics) SetupMonitoring(router *mux.Router) {
	prometheus.MustRegister(p.Hits)
	router.Handle("/metrics", promhttp.Handler())
}
