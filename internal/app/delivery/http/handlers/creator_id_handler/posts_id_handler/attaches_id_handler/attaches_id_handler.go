package attaches_id_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	useAttaches "patreon/internal/app/usecase/attaches"
	usePosts "patreon/internal/app/usecase/posts"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AttachesIDHandler struct {
	attachesUsecase useAttaches.Usecase
	bh.BaseHandler
}

func NewAttachesIDHandler(
	log *logrus.Logger,
	ucAttaches useAttaches.Usecase,
	ucPosts usePosts.Usecase,
	sClient session_client.AuthCheckerClient) *AttachesIDHandler {
	h := &AttachesIDHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		attachesUsecase: ucAttaches,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)

	h.AddMiddleware(middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost,
		middleware.NewAttachesMiddleware(log, ucAttaches).CheckCorrectAttach)

	h.AddMethod(http.MethodGet, h.GET)

	h.AddMethod(http.MethodDelete, h.DELETE,
		sessionMiddleware.CheckFunc, csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc)
	return h
}

// GET attaches
// @Summary get current attaches
// @tags attaches
// @Description get current attaches from current creator
// @Produce json
// @Success 200 {object} http_models.ResponseAttach
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 404 {object} http_models.ErrResponse "attach with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} http_models.ErrResponse "this post not belongs this creators"
// @Router /creators/{:creator_id}/posts/{:post_id}/{:attach_id} [GET]
func (h *AttachesIDHandler) GET(w http.ResponseWriter, r *http.Request) {
	var attachId int64
	var ok bool
	if attachId, ok = h.GetInt64FromParam(w, r, "attach_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	attach, err := h.attachesUsecase.GetAttach(attachId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondPost := http_models.ToResponseAttach(*attach)

	h.Log(r).Debugf("get attach with id %d", attachId)
	h.Respond(w, r, http.StatusOK, respondPost)
}

// DELETE attaches
// @Summary delete current attaches
// @tags attaches
// @Description delete current attaches from current creator
// @Produce json
// @Success 200
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/{:attach_id} [DELETE]
func (h *AttachesIDHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	var attachId int64
	var ok bool
	if attachId, ok = h.GetInt64FromParam(w, r, "attach_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	if err := h.attachesUsecase.Delete(attachId); err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}

	w.WriteHeader(http.StatusOK)
}
