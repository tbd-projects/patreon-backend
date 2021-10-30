package middleware

import (
	"net/http"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	hf "patreon/internal/app/delivery/http/handlers/base_handler/handler_interfaces"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/utilits"

	"github.com/sirupsen/logrus"
)

type CsrfMiddleware struct {
	log     utilits.LogObject
	usecase usecase_csrf.Usecase
}

func NewCsrfMiddleware(log *logrus.Logger, uc usecase_csrf.Usecase) *CsrfMiddleware {
	return &CsrfMiddleware{
		log:     utilits.NewLogObject(log),
		usecase: uc,
	}
}

// CheckCsrfTokenFunc Errors
//		Status 403 InvalidToken
//		Status 500 handler_errors.InternalError
func (mw *CsrfMiddleware) CheckCsrfTokenFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := utilits.Responder{LogObject: mw.log}

		csrfTokenFromHeader := r.Header.Get("X-Csrf-Token")
		sessionId, okSession := r.Context().Value("session_id").(string)
		userId, okUser := r.Context().Value("user_id").(int64)
		if !okSession || !okUser {
			mw.log.Log(r).Info("can not get user_id form context")
			mw.log.Log(r).Infof("user_id = %v sessionID = %v", userId, sessionId)
			respond.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
			return
		}
		if len(csrfTokenFromHeader) == 0 {
			mw.log.Log(r).Infof("csrf token from header is empty")
			mw.log.Log(r).Infof("userId: %v sessionId: %v", userId, sessionId)
			respond.Error(w, r, http.StatusForbidden, InvalidToken)
			return
		}
		if err := mw.usecase.Check(sessionId, userId, csrfTokenFromHeader); err != nil {
			mw.log.Log(r).Infof("CSRF Middleware. CSRF token expired or not valid. err: %v", err)
			respond.Error(w, r, http.StatusForbidden, InvalidToken)
			return
		}
		next(w, r)
	}
}

func (mw *CsrfMiddleware) CheckCsrfToken(next http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCsrfTokenFunc(next.ServeHTTP))
}
