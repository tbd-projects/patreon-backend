package creator_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	db_models "patreon/internal/app/models"
	usecase_creator "patreon/internal/app/usecase/creator"
	usecase_user "patreon/internal/app/usecase/user"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	"patreon/internal/microservices/auth/sessions/middleware"

	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type CreatorHandler struct {
	sessionClient  session_client.AuthCheckerClient
	creatorUsecase usecase_creator.Usecase
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewCreatorHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient,
	ucCreator usecase_creator.Usecase, ucUser usecase_user.Usecase) *CreatorHandler {
	h := &CreatorHandler{
		BaseHandler:    *bh.NewBaseHandler(log),
		creatorUsecase: ucCreator,
		sessionClient:  sManager,
		userUsecase:    ucUser,
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMethod(http.MethodPost, h.POST, middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// GET Creators
// @Summary get list of Creators
// @Description get list of creators which register on service
// @Produce json
// @tags creators
// @Success 201 {array} http_models.ResponseCreator
// @Failure 403 {object} http_models.ErrResponse "csrf token is invalid, get new token"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation"
// @Router /creators [GET]
func (h *CreatorHandler) GET(w http.ResponseWriter, r *http.Request) {
	creators, err := h.creatorUsecase.GetCreators()
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondCreators := make([]http_models.ResponseCreator, len(creators))
	for i, cr := range creators {
		respondCreators[i] = http_models.ToResponseCreator(cr)
	}

	h.Log(r).Debugf("get creators %v", respondCreators)
	h.Respond(w, r, http.StatusOK, respondCreators)
}

// POST Create Creator
// @Summary create creator
// @Description create creator with id from path, and respond created creator
// @Param creator body http_models.RequestCreator true "Request body for creators"
// @Produce json
// @tags creators
// @Success 201 {object} http_models.IdResponse
// @Failure 409 {object} http_models.ErrResponse "creator already exist"
// @Failure 404 {object} http_models.ErrResponse "user with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid creator nickname", "invalid creator category-description", "invalid creator category", "invalid body in request"
// @Failure 403 {object} http_models.ErrResponse "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators [POST]
func (h *CreatorHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestCreator{}
	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u, err := h.userUsecase.GetProfile(userID.(int64))
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	cr := &db_models.Creator{
		ID:          u.ID,
		Nickname:    u.Nickname,
		Category:    req.Category,
		Description: req.Description,
	}

	creatorId, err := h.creatorUsecase.Create(cr)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	h.Respond(w, r, http.StatusCreated, http_models.IdResponse{ID: creatorId})
}
