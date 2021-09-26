package middleware

import (
	"context"
	"net/http"
	"patreon/internal/app/sessions"

	"github.com/sirupsen/logrus"
)

type SessionMiddleware struct {
	SessionManager sessions.SessionsManager
	log            *logrus.Logger
}

func NewSessionMiddleware(sessionManager sessions.SessionsManager, log *logrus.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		SessionManager: sessionManager,
		log:            log,
	}
}
func (m *SessionMiddleware) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			m.log.Warnf("in parsing cookie: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		uniqID := sessionID.Value
		if res, err := m.SessionManager.Check(uniqID); err != nil {
			m.log.Warnf("Error in checking session: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			m.log.Infof("Get session for user: %d", res.UserID)
			r = r.WithContext(context.WithValue(r.Context(), "user_id", res.UserID)) //nolint
			r = r.WithContext(context.WithValue(r.Context(), "uniq_id", res.UniqID)) //nolint
		}
		next.ServeHTTP(w, r)
	})
}

func (m *SessionMiddleware) CheckNotAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			m.log.Debug("User not Authorized")
			next.ServeHTTP(w, r)
			return
		}

		uniqID := sessionID.Value
		if res, err := m.SessionManager.Check(uniqID); err != nil {
			m.log.Debug("User not Authorized")
			next.ServeHTTP(w, r)
			return
		} else {
			m.log.Warnf("UserAuthorized: %d", res.UserID)
		}
		w.WriteHeader(http.StatusTeapot)
	})
}
func (m *SessionMiddleware) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			m.log.Info("options request")
			w.WriteHeader(http.StatusNoContent)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
