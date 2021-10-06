package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/handlers/base_handler"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/models"

	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	dataStorage app.DataStorage
	bh.BaseHandler
}

func NewProfileHandler(log *logrus.Logger, storage app.DataStorage) *ProfileHandler {
	h := &ProfileHandler{
		BaseHandler: *bh.NewBaseHandler(log),
		dataStorage: storage,
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.dataStorage.SessionManager(), h.Log()).Check)
	return h
}

// Profile
// @Summary get information from user for profile
// @Description get nickname and avatar for profile
// @Accept  json
// @Produce json
// @Success 201 {object} models.ProfileResponse "Successfully get profile"
// @Failure 401 "User are not authorized"
// @Failure 503 {object} models.BaseResponse "Not found user"
// @Failure 500 {object} models.BaseResponse "Error context"
// @Router /profile [GET]
func (h *ProfileHandler) GET(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log().Error(err)
		}
	}(r.Body)

	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log().Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ContextError)
		return
	}

	u, err := h.dataStorage.Store().User().FindByID(userID.(int64))
	if err != nil {
		h.Log().Errorf("get: %s err:%s can not get user from db", u, err)
		h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}

	h.Log().Debugf("get profile %s", u)
	h.Respond(w, r, http.StatusOK, models.Profile{Nickname: u.Nickname, Avatar: u.Avatar})
}
