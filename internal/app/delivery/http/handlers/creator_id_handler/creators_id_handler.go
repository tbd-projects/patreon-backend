package creator_id_handler

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/sessions"
	usecase_creator "patreon/internal/app/usecase/creator"
	usecase_user "patreon/internal/app/usecase/user"

	"github.com/gorilla/mux"
)

type CreatorIdHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	creatorUsecase usecase_creator.Usecase
	base_handler.BaseHandler
}

func NewCreatorIdHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucUser usecase_user.Usecase, ucCreator usecase_creator.Usecase) *CreatorIdHandler {
	h := &CreatorIdHandler{
		BaseHandler:    *base_handler.NewBaseHandler(log, router, cors),
		sessionManager: sManager,
		userUsecase:    ucUser,
		creatorUsecase: ucCreator,
	}
	h.AddMethod(http.MethodGet, h.GET)
	return h
}

// GET Creator
// @Summary get creator
// @Description get creator with id from path
// @Produce json
// @Param id path int true "Get creator with id"
// @Success 200 {object} models.Creator
// @Failure 404 {object} models.ErrResponse "user with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Router /creators/{creator_id:} [GET]
func (s *CreatorIdHandler) GET(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.Log(r).Error(err)
		}
	}(r.Body)

	creatorId, ok := s.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		s.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		s.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	creator, err := s.creatorUsecase.GetCreator(creatorId)
	if err != nil {
		s.UsecaseError(w, r, err, codesByErrors)

		return
	}

	s.Log(r).Debugf("get creator %v with id %v", creator, creatorId)
	s.Respond(w, r, http.StatusOK, creator)

}
