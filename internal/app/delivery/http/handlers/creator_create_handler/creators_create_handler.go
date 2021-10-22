package creator_create_handler

import (
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_creator "patreon/internal/app/usecase/creator"
	usecase_user "patreon/internal/app/usecase/user"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type CreatorCreateHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	creatorUsecase usecase_creator.Usecase
	base_handler.BaseHandler
}

func NewCreatorCreateHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucUser usecase_user.Usecase, ucCreator usecase_creator.Usecase) *CreatorCreateHandler {
	h := &CreatorCreateHandler{
		BaseHandler:    *base_handler.NewBaseHandler(log, router, cors),
		sessionManager: sManager,
		userUsecase:    ucUser,
		creatorUsecase: ucCreator,
	}
	h.AddMethod(http.MethodGet, h.GET)

	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check)
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
// @Router /creators/{:id} [GET]
func (s *CreatorCreateHandler) GET(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.Log(r).Error(err)
		}
	}(r.Body)

	vars := mux.Vars(r)
	id, ok := vars["id"]
	s.Log(r).Info("in /creators/id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if len(vars) > 1 || !ok || err != nil {
		s.Log(r).Info(vars)
		s.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	creator, err := s.creatorUsecase.GetCreator(idInt)
	if err != nil {
		s.UsecaseError(w, r, err, codesByErrors)

		return
	}

	s.Log(r).Debugf("get creator %v with id %v", creator, id)
	s.Respond(w, r, http.StatusOK, creator)

}
