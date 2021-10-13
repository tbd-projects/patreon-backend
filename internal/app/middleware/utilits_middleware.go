package middleware

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"time"
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
		var logger *logrus.Entry
		logger.Logger = mw.log
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
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		})
		r = r.WithContext(context.WithValue(r.Context(), "logger", upgradeLogger)) //nolint
		handler.ServeHTTP(w, r)
	})
}
