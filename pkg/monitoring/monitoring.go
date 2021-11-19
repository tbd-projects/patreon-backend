package monitoring

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type Monitoring interface {
	SetupMonitoring(router *mux.Router)
	GetHits() *prometheus.CounterVec
}
