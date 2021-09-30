package handlers

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/handlers/urls"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"time"

	gh "patreon/internal/app/handlers/general_handlers"

	"github.com/sirupsen/logrus"
)

type LogoutHandler struct {
	authMiddleware middleware.SessionMiddleware
	SessionManager sessions.SessionsManager
	gh.RespondHandler
	withHideMethod
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{
		RespondHandler: gh.RespondHandler{},
		withHideMethod: withHideMethod{gh.NewBaseHandler(logrus.New(), urls.Logout)},
	}
}

func (h *LogoutHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.Log())
}

func (h *LogoutHandler) Join(router *mux.Router) {
	h.baseHandler.AddMethod(gh.GET, h.ServeHTTP)
	h.baseHandler.AddMethod(gh.OPTIONAL, h.ServeHTTP)
	h.baseHandler.AddMiddleware(h.authMiddleware.Check)
	h.baseHandler.Join(router)
}

// Profile
// @Summary logout user
// @Description logout user
// @Accept  json
// @Produce json
// @Success 201 {object} models.BaseResponse "Successfully logout"
// @Failure 500 {object} models.BaseResponse "Error logout session"
// @Failure 401 "User not are authorized"
// @Router /logout [GET]
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log().Fatal(err)
		}
	}(r.Body)
	uniqID := r.Context().Value("uniq_id")
	if uniqID == nil {
		h.Log().Error("can not get uniq_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ContextError)
		return
	}

	h.Log().Debugf("Logout session: %s", uniqID)

	err := h.SessionManager.Delete(uniqID.(string))
	if err != nil {
		h.Log().Errorf("can not delete session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.DeleteCookieFail)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   uniqID.(string),
		Expires: time.Now().AddDate(0, 0, -1),
	}
	http.SetCookie(w, cookie)
	h.Respond(w, r, http.StatusOK, "successfully logout")
}
