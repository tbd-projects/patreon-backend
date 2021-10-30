package upl_text_data_handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	models_db "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"
	usePostsData "patreon/internal/app/usecase/posts_data"

	"github.com/sirupsen/logrus"
)

type PostsDataUploadTextHandler struct {
	postsDataUsecase usePostsData.Usecase
	bh.BaseHandler
}

func NewPostsDataUploadTextHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	ucPostsData usePostsData.Usecase, ucPosts usePosts.Usecase,
	manager sessions.SessionsManager) *PostsDataUploadTextHandler {
	h := &PostsDataUploadTextHandler{
		BaseHandler:      *bh.NewBaseHandler(log, router, cors),
		postsDataUsecase: ucPostsData,
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
// @Param user body models.RequestText true "Request body for text"
// @Success 201 {object} models.IdResponse "id posts_data"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 422 {object} models.ErrResponse "invalid data type"
// @Failure 422 {object} models.ErrResponse "this post id not know"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Router /creators/{:creator_id}/posts/{:post_id}/text [POST]
func (h *PostsDataUploadTextHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

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

	req := &models.RequestText{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	dataId, err := h.postsDataUsecase.LoadText(&models_db.PostData{Data: req.Text, PostId: postId})
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
	}

	h.Respond(w, r, http.StatusCreated, &models.IdResponse{ID: dataId})
}