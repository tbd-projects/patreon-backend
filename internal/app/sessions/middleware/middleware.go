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
			m.log.Errorf("Error in parsing cookie: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		uniqID := sessionID.Value
		if res, err := m.SessionManager.Check(uniqID); err != nil {
			m.log.Warnf("Error in checking session: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			m.log.Infof("Get session for user: %d\n", res.UserID)
			r = r.WithContext(context.WithValue(r.Context(), "user_id", res)) //nolint
		}
		next.ServeHTTP(w, r)
	})
}

func (m *SessionMiddleware) CheckNotAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			m.log.Info("User not Authorized")
			next.ServeHTTP(w, r)
			return
		}

		uniqID := sessionID.Value
		if res, err := m.SessionManager.Check(uniqID); err != nil {
			m.log.Info("User not Authorized")
			next.ServeHTTP(w, r)
			return
		} else {
			m.log.Errorf("UserAuthorized: %d\n", res.UserID)
		}
		w.WriteHeader(http.StatusTeapot)
	})
}