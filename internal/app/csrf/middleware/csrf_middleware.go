package middleware

import (
	"net/http"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"

	"github.com/sirupsen/logrus"
)

type CsrfMiddleware struct {
	log     *logrus.Logger
	usecase usecase_csrf.Usecase
}

func NewCsrfMiddleware(log *logrus.Logger, uc usecase_csrf.Usecase) *CsrfMiddleware {
	return &CsrfMiddleware{
		log:     log,
		usecase: uc,
	}
}

func (mw *CsrfMiddleware) CheckCsrfToken(next bh.HandlerFunc) bh.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfTokenFromHeader := r.Header.Get("X-Csrf-Token")
		sessionId, okSession := r.Context().Value("session_id").(string)
		userId, okUser := r.Context().Value("user_id").(int64)
		if !okSession || !okUser {
			mw.log.Infof("can not get userId or sessionId from context")
			mw.log.Infof("userId: %v sessionId: %v", userId, sessionId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(csrfTokenFromHeader) == 0 {
			mw.log.Infof("csrf token from header is empty")
			mw.log.Infof("userId: %v sessionId: %v", userId, sessionId)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if err := mw.usecase.Check(sessionId, userId, csrfTokenFromHeader); err != nil {
			mw.log.Infof("CSRF Middleware. CSRF token expired or not valid. err: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
