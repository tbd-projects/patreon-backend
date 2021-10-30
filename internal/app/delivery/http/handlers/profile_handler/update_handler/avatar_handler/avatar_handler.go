package avatar_handler

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"patreon/internal/app"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_user "patreon/internal/app/usecase/user"

	"github.com/sirupsen/logrus"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 4 // 4MB

type UpdateAvatarHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewUpdateAvatarHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	sManager sessions.SessionsManager, ucUser usecase_user.Usecase) *UpdateAvatarHandler {
	h := &UpdateAvatarHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
	}
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken)
	h.AddMethod(http.MethodPut, h.PUT)
	return h
}

// PUT AvatarChange
// @Summary set new user avatar
// @Accept  image/png, image/jpeg, image/jpg
// @Param file formData file true "Avatar file with ext jpeg/png"
// @Success 200 "successfully upload avatar"
// @Failure 400 {object} models.ErrResponse "size of file very big"
// @Failure 400 {object} models.ErrResponse "invalid form field name"
// @Failure 400 {object} models.ErrResponse "invalid avatar extension"
// @Failure 403 "csrf token is invalid, get new token"
// @Failure 500 {object} models.ErrResponse "server error"
// @Router /user/update/avatar [PUT]
func (h *UpdateAvatarHandler) PUT(w http.ResponseWriter, r *http.Request) {
	file, filename, code, err := h.GerFilesFromRequest(w, r, bh.MAX_UPLOAD_SIZE,
		"avatar", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
	}

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			Err:         handler_errors.InternalError,
			ExternalErr: errors.New("context parse userId error"),
		})
		return
	}

	err = h.userUsecase.UpdateAvatar(file, filename, userId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
	}
	w.WriteHeader(http.StatusOK)
}
