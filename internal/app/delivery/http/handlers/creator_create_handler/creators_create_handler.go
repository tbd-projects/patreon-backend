package creator_create_handler

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	db_models "patreon/internal/app/models"
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

func NewCreatorCreateHandler(log *logrus.Logger, sManager sessions.SessionsManager,
	ucUser usecase_user.Usecase, ucCreator usecase_creator.Usecase) *CreatorCreateHandler {
	h := &CreatorCreateHandler{
		BaseHandler:    *base_handler.NewBaseHandler(log),
		sessionManager: sManager,
		userUsecase:    ucUser,
		creatorUsecase: ucCreator,
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMethod(http.MethodPost, h.POST)

	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check)
	return h
}

// POST Create Creator
// @Summary create creator
// @Description create creator with id from path, and respond created creator
// @Produce json
// @Param creator body models.RequestCreator true "Request body for create"
// @Success 200 {object} models.Creator "Create user successfully"
// @Failure 400 {object} models.BaseResponse "Invalid request query"
// @Failure 404 {object} models.BaseResponse "User with id not found"
// @Failure 409 {object} models.BaseResponse "Creator already exist"
// @Failure 422 {object} models.BaseResponse "Invalid request body"
// @Failure 503 {object} models.BaseResponse "Internal error"
// @Router /creators/{:id} [POST]
func (s *CreatorCreateHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.Log(r).Error(err)
		}
	}(r.Body)

	vars := mux.Vars(r)
	id, ok := vars["id"]
	s.Log(r).Info("in /creators/id")
	idInt, err := strconv.Atoi(id)
	if len(vars) > 1 || !ok || err != nil {
		s.Log(r).Info(vars)
		s.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	req := &models.RequestCreator{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		s.Log(r).Warnf("can not parse request %s", err)
		s.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	u, err := s.userUsecase.GetProfile(int64(idInt))
	if err != nil {
		s.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}
	cr := &db_models.Creator{
		ID:          u.ID,
		Nickname:    u.Nickname,
		Category:    req.Category,
		Description: req.Description,
	}
	creatorId, err := s.creatorUsecase.Create(cr)
	if err != nil {
		s.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	s.Respond(w, r, http.StatusOK, creatorId)
}

// GET Creator
// @Summary get creator
// @Description get creator with id from path
// @Produce json
// @Param id path int true "Get creator with id"
// @Success 200 {object} models.Creator "Get user successfully"
// @Failure 503 {object} models.BaseResponse "Internal error"
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
		s.UsecaseError(w, r, err, codesByErrorsGET)

		return
	}

	s.Log(r).Debugf("get creator %v with id %v", creator, id)
	s.Respond(w, r, http.StatusOK, creator)

}
