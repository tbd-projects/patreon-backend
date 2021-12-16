package prometheus_monitoring

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
			Name: serviceName + "_success_hits",
			Help: "Count success responses from service",
		}, []string{"status", "path", "method"}),
		HitsErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: serviceName + "_errors_hits",
			Help: "Count errors response from service",
		}, []string{"status", "path", "method"}),
		ExecutionTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: serviceName + "_durations",
			Help: "Duration execution of request",
		}, []string{"status", "path", "method"}),
		TotalHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: serviceName + "_total_hits",
		}),
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
func (pm *PrometheusMetrics) GetSuccessHits() *prometheus.CounterVec {
	return pm.HitsSuccess
}
func (pm *PrometheusMetrics) GetErrorsHits() *prometheus.CounterVec {
	return pm.HitsErrors
}
func (pm *PrometheusMetrics) GetRequestCounter() prometheus.Counter {
	return pm.TotalHits
}
func (pm *PrometheusMetrics) GetExecution() *prometheus.HistogramVec {
	return pm.ExecutionTime
}

func CreateNewMonitoringRouter(host string) {
	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf("%s:9091", host), router); err != nil {
		log.Fatalln(err)
	}
}
