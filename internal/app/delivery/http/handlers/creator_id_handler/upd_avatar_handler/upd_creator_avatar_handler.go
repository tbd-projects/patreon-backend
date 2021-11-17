package upd_avatar_creator_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	middlewareSes "patreon/internal/app/sessions/middleware"
	usecase_creator "patreon/internal/app/usecase/creator"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UpdateAvatarCreatorHandler struct {
	sessionClient  session_client.AuthCheckerClient
	creatorUsecase usecase_creator.Usecase
	bh.BaseHandler
}

func NewUpdateAvatarHandler(log *logrus.Logger,
	sessionClient session_client.AuthCheckerClient, creatorUsecase usecase_creator.Usecase) *UpdateAvatarCreatorHandler {
	h := &UpdateAvatarCreatorHandler{
		sessionClient:  sessionClient,
		creatorUsecase: creatorUsecase,
		BaseHandler:    *bh.NewBaseHandler(log),
	}
	h.AddMiddleware(middlewareSes.NewSessionMiddleware(h.sessionClient, log).Check,
		middleware.NewCreatorsMiddleware(log).CheckAllowUser)
	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// PUT AvatarChange
// @Summary set new creator avatar
// @Accept  image/png, image/jpeg, image/jpg
// @Param avatar formData file true "Avatar file with ext jpeg/png"
// @Success 200 "successfully upload avatar"
// @Failure 400 {object} models.ErrResponse "size of file very big", "please upload a JPEG, JPG or PNG files", "invalid form field name"
// @Failure 403 {object} models.ErrResponse "csrf token is invalid, get new token"
// @Failure 422 {object} models.ErrResponse "this creator id not know"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Router /creators/{creator_id:}/update/avatar [PUT]
func (h *UpdateAvatarCreatorHandler) PUT(w http.ResponseWriter, r *http.Request) {
	file, filename, code, err := h.GerFilesFromRequest(w, r, bh.MAX_UPLOAD_SIZE,
		"avatar", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	creatorId, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	err = h.creatorUsecase.UpdateAvatar(file, filename, creatorId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
