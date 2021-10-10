package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/sessions/middleware"
	"time"

	"github.com/sirupsen/logrus"
)

type LogoutHandler struct {
	dataStorage app.DataStorage
	bh.BaseHandler
}

func NewLogoutHandler(log *logrus.Logger, storage app.DataStorage) *LogoutHandler {
	h := &LogoutHandler{
		BaseHandler: *bh.NewBaseHandler(log),
		dataStorage: storage,
	}
	h.AddMethod(http.MethodPost, h.POST)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.dataStorage.SessionManager(), h.Log()).Check)
	return h
}

// Profile
// @Summary logout user
// @Description logout user
// @Accept  json
// @Produce json
// @Success 201 {object} models.BaseResponse "Successfully logout"
// @Failure 500 {object} models.BaseResponse "Error logout session"
// @Failure 401 "User not are authorized"
// @Router /logout [POST]
func (h *LogoutHandler) POST(w http.ResponseWriter, r *http.Request) {
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

	err := h.dataStorage.SessionManager().Delete(uniqID.(string))
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
