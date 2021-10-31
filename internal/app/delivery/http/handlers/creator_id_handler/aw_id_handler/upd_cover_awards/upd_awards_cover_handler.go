package upd_cover_awards_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	"patreon/internal/app/sessions"
	middlewareSes "patreon/internal/app/sessions/middleware"
	useAwards "patreon/internal/app/usecase/awards"
)

type UpdateCoverAwardsHandler struct {
	sessionManager sessions.SessionsManager
	awardsUsecase  useAwards.Usecase
	bh.BaseHandler
}

func NewUpdateCoverAwardsHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	sManager sessions.SessionsManager, awardsUsecase useAwards.Usecase) *UpdateCoverAwardsHandler {
	h := &UpdateCoverAwardsHandler{
		sessionManager: sManager,
		awardsUsecase:  awardsUsecase,
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
	}
	h.AddMiddleware(middlewareSes.NewSessionMiddleware(h.sessionManager, log).Check,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken,
		middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewAwardsMiddleware(log, awardsUsecase).CheckCorrectAward)
	h.AddMethod(http.MethodPut, h.PUT)
	return h
}

// PUT CoverChange
// @Summary set new awards cover
// @Accept  image/png, image/jpeg, image/jpg
// @Param cover formData file true "Cover file with ext jpeg/png"
// @Success 200 "successfully upload avatar"
// @Failure 400 {object} models.ErrResponse "size of file very big"
// @Failure 400 {object} models.ErrResponse "invalid form field name"
// @Failure 400 {object} models.ErrResponse "please upload a JPEG, JPG or PNG files"
// @Failure 403 "csrf token is invalid, get new token"
// @Failure 422 {object} models.ErrResponse "this creator id not know"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Router /creators/{:creator_id}/awards/{:award_id}/cover [PUT]
func (h *UpdateCoverAwardsHandler) PUT(w http.ResponseWriter, r *http.Request) {
	file, filename, code, err := h.GerFilesFromRequest(w, r, bh.MAX_UPLOAD_SIZE,
		"cover", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	awardId, ok := h.GetInt64FromParam(w, r, "award_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	err = h.awardsUsecase.UpdateCover(file, filename, awardId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
