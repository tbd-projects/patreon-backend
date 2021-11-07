package creator_id_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_creator "patreon/internal/app/usecase/creator"

	"github.com/gorilla/mux"
)

type CreatorIdHandler struct {
	sessionManager sessions.SessionsManager
	creatorUsecase usecase_creator.Usecase
	base_handler.BaseHandler
}

func NewCreatorIdHandler(log *logrus.Logger, sManager sessions.SessionsManager, ucCreator usecase_creator.Usecase) *CreatorIdHandler {
	h := &CreatorIdHandler{
		BaseHandler:    *base_handler.NewBaseHandler(log),
		sessionManager: sManager,
		creatorUsecase: ucCreator,
	}
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).AddUserId)
	h.AddMethod(http.MethodGet, h.GET)

	return h
}

// GET Creator
// @Summary get creator
// @Description get creator with id from path
// @Produce json
// @Param creator_id path int true "Get creator with id"
// @Success 200 {object} models.ResponseCreatorWithAwards
// @Failure 404 {object} models.ErrResponse "user not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Router /creators/{creator_id:} [GET]
func (s *CreatorIdHandler) GET(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		userId = usecase_creator.NoUser
	}

	creatorId, ok := s.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		s.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		s.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	creator, err := s.creatorUsecase.GetCreator(creatorId, userId)
	if err != nil {
		s.UsecaseError(w, r, err, codesByErrors)
		return
	}

	s.Log(r).Debugf("get creator %v with id %v", creator, creatorId)
	s.Respond(w, r, http.StatusOK, models.ToResponseCreatorWithAwards(*creator))
}
