package upl_text_data_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	models_db "patreon/internal/app/models"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"
	usePostsData "patreon/internal/app/usecase/posts_data"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type PostsDataUploadTextHandler struct {
	postsDataUsecase usePostsData.Usecase
	bh.BaseHandler
}

func NewPostsDataUploadTextHandler(log *logrus.Logger,
	ucPostsData usePostsData.Usecase, ucPosts usePosts.Usecase,
	sClient session_client.AuthCheckerClient) *PostsDataUploadTextHandler {
	h := &PostsDataUploadTextHandler{
		BaseHandler:      *bh.NewBaseHandler(log),
		postsDataUsecase: ucPostsData,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(sClient, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)
	h.AddMethod(http.MethodPost, h.POST)
	return h
}

// POST add text to post
// @Summary add text to post
// @Accept  json
// @Param text body models.RequestText true "Request body for text"
// @Success 201 {object} models.IdResponse "id posts_data"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} models.ErrResponse "invalid data type", "this post id not know"
// @Failure 400 {object} models.ErrResponse "invalid parameters", "invalid body in request"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/text [POST]
func (h *PostsDataUploadTextHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &models.RequestText{}

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

	dataId, err := h.postsDataUsecase.LoadText(&models_db.PostData{Data: req.Text, PostId: postId})
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Respond(w, r, http.StatusCreated, &models.IdResponse{ID: dataId})
}
