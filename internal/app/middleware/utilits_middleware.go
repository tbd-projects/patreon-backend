package middleware

import (
	"net/http"
	"patreon/internal/app/utilits"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type UtilitiesMiddleware struct {
	log utilits.LogObject
}

func NewUtilitiesMiddleware(log *logrus.Logger) UtilitiesMiddleware {
	return UtilitiesMiddleware{utilits.NewLogObject(log)}
}

func (mw *UtilitiesMiddleware) CheckPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(log *logrus.Entry, w http.ResponseWriter) {
			if err := recover(); err != nil {
				log.Errorf("detacted critical error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}(mw.log.Log(r), w)
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
		handler.ServeHTTP(w, r)
		upgradeLogger.Infof("work time: %v", time.Since(start).Milliseconds())
	})
}
