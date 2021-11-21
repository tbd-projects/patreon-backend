package prometheus_monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	HitsSuccess   *prometheus.CounterVec
	HitsErrors    *prometheus.CounterVec
	ExecutionTime *prometheus.HistogramVec
	TotalHits     prometheus.Counter
}

func NewPrometheusMetrics(serviceName string) *PrometheusMetrics {
	metrics := &PrometheusMetrics{
		HitsSuccess: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: serviceName + "_hits",
			Help: "Count success responses from service",
		}, []string{"status", "path", "method"}),
		HitsErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: serviceName + "errors",
			Help: "Count errors response from service",
		}, []string{"status", "path", "method"}),
		ExecutionTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: serviceName + "_durations",
			Help: "Duration execution of request",
		}, []string{"status", "path", "method"}),
	}

	return metrics
}
func (pm *PrometheusMetrics) SetupMonitoring() error {
	if err := prometheus.Register(pm.HitsErrors); err != nil {
		return err
	}
	if err := prometheus.Register(pm.HitsSuccess); err != nil {
		return err
	}
	if err := prometheus.Register(pm.ExecutionTime); err != nil {
		return err
	}
	if err := prometheus.Register(pm.TotalHits); err != nil {
		return err
	}
	return nil
}
