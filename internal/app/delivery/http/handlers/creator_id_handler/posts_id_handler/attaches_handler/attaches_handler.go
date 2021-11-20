package attaches_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	"patreon/internal/app/delivery/http/handlers"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	http_models "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	"patreon/internal/app/models"
	useAttaches "patreon/internal/app/usecase/attaches"
	usePosts "patreon/internal/app/usecase/posts"
	"patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/microcosm-cc/bluemonday"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AttachesHandler struct {
	attachesUsecase useAttaches.Usecase
	bh.BaseHandler
}

func NewAttachesHandler(log *logrus.Logger,
	ucAttaches useAttaches.Usecase,
	ucPosts usePosts.Usecase,
	sClient client.AuthCheckerClient) *AttachesHandler {
	h := &AttachesHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		attachesUsecase: ucAttaches,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)

	h.AddMiddleware(sessionMiddleware.Check, csrf_middleware.NewCsrfMiddleware(log,
		usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken,
		middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost)
	h.AddMethod(http.MethodPut, h.PUT)
	return h
}

// PUT attaches
// @Summary update current attaches
// @tags attaches
// @Description update current attaches from current creator
// @Produce json
// @Param attaches body http_models.RequestAttaches true "Request body for set attaches"
// @Success 200 {object} http_models.ResponseApplyAttach
// @Failure 400 {object} http_models.ErrResponse ""invalid parameters""
// @Failure 404 {object} http_models.ErrResponse ""attach with this id not found""
// @Failure 500 {object} http_models.ErrResponse ""can not do bd operation", "server error", "Not allow type, allowed type is: ...""
// @Failure 403 {object} http_models.ErrResponse ""for this user forbidden change creator", "csrf token is invalid, get new token", "this post not belongs this creators""
// @Failure 422 {object} http_models.ErrResponse ""Not allow type, allowed type is: ...", "Not valid attach id", "Not allow status, allowed status is: ...""
// @Router /creators/{:creator_id}/posts/{:post_id}/attaches [PUT]
func (h *AttachesHandler) PUT(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestAttaches{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	var postId int64
	var ok bool
	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	if err = req.Validate(); err != nil {
		h.Log(r).Warnf("invalid attach in request %v", err)
		h.Error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	var newAttaches []models.Attach
	var updateAtaches []models.Attach

	for level, attach := range req.Attaches {
		model_attach := models.Attach{Id: attach.Id, Type: attach.Type,
			Value: attach.Value, Level: int64(level + 1)}
		if attach.Type == models.Text {
			if attach.Status == handlers.AddStatus {
				newAttaches = append(newAttaches, model_attach)
			} else {
				updateAtaches = append(updateAtaches, model_attach)
			}
		} else {
			updateAtaches = append(updateAtaches, model_attach)
		}
	}

	res, err := h.attachesUsecase.UpdateAttach(postId, newAttaches, updateAtaches)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	h.Log(r).Debugf("apply attaches for post with id %d", postId)
	h.Respond(w, r, http.StatusOK, http_models.ResponseApplyAttach{IDs: res})
}
