package profile_handler

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	useUser "patreon/internal/app/usecase/user"
)

type ProfileHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    useUser.Usecase
	bh.BaseHandler
}

func NewProfileHandler(log *logrus.Logger, sManager sessions.SessionsManager, ucUser useUser.Usecase) *ProfileHandler {
	h := &ProfileHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, h.Log()).Check)
	return h
}

// Profile
// @Summary get information from user for user
// @Description get nickname and avatar for user
// @Accept  json
// @Produce json
// @Success 201 {object} models.ProfileResponse "Successfully get user"
// @Failure 401 "User are not authorized"
// @Failure 503 {object} models.BaseResponse "Not found user"
// @Failure 500 {object} models.BaseResponse "Error context"
// @Router /user [GET]
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

	u, err := h.userUsecase.GetProfile(userID.(int64))
	if err != nil {
		//h.Log().Errorf("get: %s err:%s can not get user from db", u, err)
		h.UsecaseError(w, r, err, codeByError)
		//h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}

	h.Log().Debugf("get user %s", u)
	h.Respond(w, r, http.StatusOK, models.Profile{Nickname: u.Nickname, Avatar: u.Avatar})
}
