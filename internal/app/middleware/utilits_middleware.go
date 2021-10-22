package middleware

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type UtilitiesMiddleware struct {
	log *logrus.Logger
}

func NewUtilitiesMiddleware(log *logrus.Logger) UtilitiesMiddleware {
	return UtilitiesMiddleware{log: log}
}

func (mw *UtilitiesMiddleware) CheckPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxLogger := r.Context().Value("logger")
		logger := mw.log.WithField("base_log with url:", r.URL)
		if ctxLogger != nil {
			if log, ok := ctxLogger.(*logrus.Entry); ok {
				logger = log
			}
		}

		defer func(log *logrus.Entry) {
			if err := recover(); err != nil {
				log.Errorf("detacted critical error: %v", err)
			}
		}(logger)
		handler.ServeHTTP(w, r)
	})
}

func (mw *UtilitiesMiddleware) UpgradeLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		upgradeLogger := mw.log.WithFields(logrus.Fields{
			"urls":        r.URL,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start).Milliseconds(),
			"req_id":      uuid.NewV4(),
		})
		r = r.WithContext(context.WithValue(r.Context(), "logger", upgradeLogger)) //nolint
		upgradeLogger.Info("Log was upgraded")
		handler.ServeHTTP(w, r)
	})
}
