package middleware

import (
	"net/http"
	"patreon/internal/app/utilits"
	"patreon/pkg/monitoring"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/urfave/negroni"

	uuid "github.com/satori/go.uuid"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type UtilitiesMiddleware struct {
	log     utilits.LogObject
	metrics monitoring.Monitoring
}

func NewUtilitiesMiddleware(log *logrus.Logger, metrics monitoring.Monitoring) UtilitiesMiddleware {
	return UtilitiesMiddleware{
		log:     utilits.NewLogObject(log),
		metrics: metrics,
	}
}

func (mw *UtilitiesMiddleware) CheckPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(log *logrus.Entry, metrics monitoring.Monitoring, w http.ResponseWriter) {
			if err := recover(); err != nil {
				responseErr := http.StatusInternalServerError
				metrics.GetErrorsHits().WithLabelValues(
					strconv.Itoa(responseErr),
					r.URL.String(),
					r.Method,
				)
				metrics.GetRequestCounter().Inc()

				log.Errorf("detacted critical error: %v, with stack: %s", err, debug.Stack())
				w.WriteHeader(responseErr)
			}
		}(mw.log.Log(r), mw.metrics, w)
		handler.ServeHTTP(w, r)
	})
}

func (mw *UtilitiesMiddleware) UpgradeLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		upgradeLogger := mw.log.BaseLog().WithFields(logrus.Fields{
			"urls":        r.URL,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start).Milliseconds(),
			"req_id":      uuid.NewV4(),
		})

		r = r.WithContext(context.WithValue(r.Context(), "logger", upgradeLogger))
		upgradeLogger.Info("Log was upgraded")

		wrappedWriter := negroni.NewResponseWriter(w)
		handler.ServeHTTP(wrappedWriter, r)

		statusCode := wrappedWriter.Status()

		executeTime := time.Since(start).Milliseconds()
		upgradeLogger.Infof("work time [ms]: %v", executeTime)

		mw.metrics.GetRequestCounter().Inc()

		if statusCode < 300 {
			mw.metrics.GetSuccessHits().WithLabelValues(
				strconv.Itoa(statusCode),
				r.URL.String(),
				r.Method,
			).Add(1)
		} else {
			mw.metrics.GetErrorsHits().WithLabelValues(
				strconv.Itoa(statusCode),
				r.URL.String(),
				r.Method,
			).Add(1)
		}
	})
}
