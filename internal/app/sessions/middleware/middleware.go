package middleware

import (
	"net/http"
	"patreon/internal/app/sessions/sessions_manager"

	"github.com/sirupsen/logrus"
)

type SessionMiddleware struct {
	SessionManager sessions_manager.SessionManager
	log            *logrus.Logger
}

func NewSessionMiddleware(sessionManager sessions_manager.SessionManager, log *logrus.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		SessionManager: sessionManager,
		log:            log,
	}
}
func (m *SessionMiddleware) Check(next http.HandlerFunc) http.HandlerFunc {
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
		}
		next.ServeHTTP(w, r)
	})
}
