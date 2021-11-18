package logout_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"time"

	"github.com/sirupsen/logrus"
)

type LogoutHandler struct {
	sessionManager sessions.SessionsManager
	bh.BaseHandler
}

func NewLogoutHandler(log *logrus.Logger,
	sManager sessions.SessionsManager) *LogoutHandler {
	h := &LogoutHandler{
		BaseHandler:    *bh.NewBaseHandler(log),
		sessionManager: sManager,
	}
	h.AddMethod(http.MethodPost, h.POST,
		middleware.NewSessionMiddleware(h.sessionManager, log).CheckFunc,
	)

	return h
}

// POST Logout
// @Summary logout user
// @Description logout user
// @Accept  json
// @Produce json
// @Success 201 "Successfully logout"
// @Failure 500 {object} http_models.ErrResponse "server error"
// @Failure 401 "User not are authorized"
// @Router /logout [POST]
func (h *LogoutHandler) POST(w http.ResponseWriter, r *http.Request) {
	uniqID := r.Context().Value("session_id")
	if uniqID == nil {
		h.Log(r).Error("can not get session_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	h.Log(r).Debugf("Logout session: %s", uniqID)

	err := h.sessionManager.Delete(uniqID.(string))
	if err != nil {
		h.Log(r).Errorf("can not delete session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.DeleteCookieFail)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   uniqID.(string),
		Expires: time.Now().AddDate(0, 0, -1),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
