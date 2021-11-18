package upl_text_attach_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	models_db "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	useAttaches "patreon/internal/app/usecase/attaches"
	usePosts "patreon/internal/app/usecase/posts"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type AttachesUploadTextHandler struct {
	attachesUsecase useAttaches.Usecase
	bh.BaseHandler
}

func NewAttachesUploadTextHandler(log *logrus.Logger,
	ucAttaches useAttaches.Usecase, ucPosts usePosts.Usecase,
	manager sessions.SessionsManager) *AttachesUploadTextHandler {
	h := &AttachesUploadTextHandler{
		BaseHandler:      *bh.NewBaseHandler(log),
		attachesUsecase: ucAttaches,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)
	h.AddMethod(http.MethodPost, h.POST)
	return h
}

// POST add text to post
// @Summary add text to post
// @Accept  json
// @Param text body http_models.RequestText true "Request body for text"
// @Success 201 {object} http_models.IdResponse "id attaches"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid data type", "this post id not know"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters", "invalid body in request"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/text [POST]
func (h *AttachesUploadTextHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestText{}

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

	attachId, err := h.attachesUsecase.LoadText(&models_db.PostData{Data: req.Text, PostId: postId})
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Respond(w, r, http.StatusCreated, &http_models.IdResponse{ID: attachId})
}
