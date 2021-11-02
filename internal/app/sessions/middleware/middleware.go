package middleware

import (
	"context"
	"net/http"
	hf "patreon/internal/app/delivery/http/handlers/base_handler/handler_interfaces"
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

// CheckFunc Errors:
//		Status 401 "not authorized user"
func (m *SessionMiddleware) CheckFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			m.log.Debugf("Get session for user: %d", res.UserID)
			r = r.WithContext(context.WithValue(r.Context(), "user_id", res.UserID))    //nolint
			r = r.WithContext(context.WithValue(r.Context(), "session_id", res.UniqID)) //nolint
		}
		next(w, r)
	}
}

// Check Errors:
//		Status 401 "not authorized user"
func (m *SessionMiddleware) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(m.CheckFunc(next.ServeHTTP))
}

// CheckNotAuthorized Errors:
//		Status 418 "user already authorized"
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

// AddUserIdFunc Errors:
//		Nothing return only add user_id and session_id to context
func (m *SessionMiddleware) AddUserIdFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err == nil {
			uniqID := sessionID.Value
			if res, err := m.SessionManager.Check(uniqID); err == nil {
				m.log.Debugf("Get session for user: %d", res.UserID)
				r = r.WithContext(context.WithValue(r.Context(), "user_id", res.UserID))
				r = r.WithContext(context.WithValue(r.Context(), "session_id", res.UniqID))
			}
		}
		next(w, r)
	}
}

// AddUserId Errors:
//		Nothing return only add user_id and session_id to context
func (m *SessionMiddleware) AddUserId(next http.Handler) http.Handler {
	return http.HandlerFunc(m.AddUserIdFunc(next.ServeHTTP))
}
