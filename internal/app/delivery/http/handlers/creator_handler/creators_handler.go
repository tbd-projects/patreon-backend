package creator_handler

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	db_models "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	middleSes "patreon/internal/app/sessions/middleware"
	usecase_creator "patreon/internal/app/usecase/creator"
	usecase_user "patreon/internal/app/usecase/user"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type CreatorHandler struct {
	sessionManager sessions.SessionsManager
	creatorUsecase usecase_creator.Usecase
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewCreatorHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucCreator usecase_creator.Usecase, ucUser usecase_user.Usecase) *CreatorHandler {
	h := &CreatorHandler{
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
		creatorUsecase: ucCreator,
		sessionManager: sManager,
		userUsecase:    ucUser,
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMethod(http.MethodPost, h.POST, middleSes.NewSessionMiddleware(h.sessionManager, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// GET Creators
// @Summary get list of Creators
// @Description get list of creators which register on service
// @Produce json
// @Success 201 {array} models.ResponseCreator
// @Failure 403 "csrf token is invalid, get new token"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Router /creators [GET]
func (h *CreatorHandler) GET(w http.ResponseWriter, r *http.Request) {
	creators, err := h.creatorUsecase.GetCreators()
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondCreators := make([]models.ResponseCreator, len(creators))
	for i, cr := range creators {
		respondCreators[i] = models.ToResponseCreator(cr)
	}

	h.Log(r).Debugf("get creators %v", respondCreators)
	h.Respond(w, r, http.StatusOK, respondCreators)
}

// POST Create Creator
// @Summary create creator
// @Description create creator with id from path, and respond created creator
// @Param user body models.RequestCreator true "Request body for creators"
// @Produce json
// @Success 201 {object} models.Creator
// @Failure 422 {object} models.ErrResponse "invalid body in request"
// @Failure 500 {object} models.ErrResponse "server error
// @Failure 409 {object} models.ErrResponse "creator already exist"
// @Failure 404 {object} models.ErrResponse "user with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 422 {object} models.ErrResponse "invalid creator category"
// @Failure 422 {object} models.ErrResponse "invalid creator nickname"
// @Failure 422 {object} models.ErrResponse "invalid creator category-description"
// @Failure 401 "User are not authorized"
// @Router /creators [POST]
func (s *CreatorHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.Log(r).Error(err)
		}
	}(r.Body)

	req := &models.RequestCreator{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(req); err != nil {
		s.Log(r).Warnf("can not parse request %s", err)
		s.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	userID := r.Context().Value("user_id")
	if userID == nil {
		s.Log(r).Error("can not get user_id from context")
		s.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u, err := s.userUsecase.GetProfile(userID.(int64))
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

	s.Respond(w, r, http.StatusCreated, creatorId)
}
