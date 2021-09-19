package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/sessions"
)

type SessionsMiddleware struct {
	sessionsManager sessions.ISessionsManager
	log             *logrus.Logger
}

func CreateSessionMiddleware(sessionManager sessions.ISessionsManager, log *logrus.Logger) *SessionsMiddleware {
	sessMiddleware := &SessionsMiddleware{
		sessionsManager: sessionManager,
		log:             log,
	}
	return sessMiddleware
}

func (sessMiddleware *SessionsMiddleware) CheckSession(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			sessMiddleware.log.Errorf("Error in parsing cookie: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		uniqID := sessionID.Value
		result, err := sessMiddleware.sessionsManager.CheckSession(uniqID)

		if err == nil {
			sessMiddleware.log.Infof("Get session for user: %d\n", result.UserID)
		} else {
			sessMiddleware.log.Warnf("Error in checking session: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
