package middleware

import (
	"net/http"
	usecase_csrf "patreon/internal/app/csrf/usecase"

	"github.com/sirupsen/logrus"
)

type CsrfMiddleware struct {
	log     *logrus.Logger
	usecase usecase_csrf.Usecase
}

func NewCsrfMiddleware(log *logrus.Logger, uc usecase_csrf.Usecase) CsrfMiddleware {
	return CsrfMiddleware{
		log:     log,
		usecase: uc,
	}
}

func (mw *CsrfMiddleware) CheckCsrfToken(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfTokenFromHeader := r.Header.Get("X-Csrf-Token")
		csrfTokenFromCookie, err := r.Cookie("csrf")
		if err != nil {
			mw.log.Infof("in parsing cookie: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if csrfTokenFromCookie.Value != csrfTokenFromHeader {
			mw.log.Infof("No CSRF Token header: %v cookie %v", csrfTokenFromHeader, csrfTokenFromCookie.Value)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		sessionId, okSession := r.Context().Value("session_id").(string)
		userId, okUser := r.Context().Value("user_id").(int64)
		if !okSession || !okUser {
			mw.log.Infof("can not get userId or sessionId from context")
			mw.log.Infof("userId: %v sessionId: %v", userId, sessionId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if mw.usecase.Check(sessionId, userId, csrfTokenFromHeader) != nil {
			mw.log.Infof("CSRF token expired or not valid. err: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
