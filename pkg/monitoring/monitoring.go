package monitoring

import "github.com/gorilla/mux"

type Monitoring interface {
	SetupMonitoring(router *mux.Router)
}
