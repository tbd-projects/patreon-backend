package posts_data_id_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"
	usePostsData "patreon/internal/app/usecase/posts_data"
)

type PostsDataIDHandler struct {
	postsDataUsecase usePostsData.Usecase
	bh.BaseHandler
}

func NewPostsDataIDHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	ucPostsData usePostsData.Usecase, ucPosts usePosts.Usecase, manager sessions.SessionsManager) *PostsDataIDHandler {
	h := &PostsDataIDHandler{
		BaseHandler:      *bh.NewBaseHandler(log, router, cors),
		postsDataUsecase: ucPostsData,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	h.AddMiddleware(middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMethod(http.MethodDelete, h.DELETE, sessionMiddleware.CheckFunc, middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc)
	return h
}

// GET posts_data
// @Summary get current posts_data
// @Description get current posts_data from current creator
// @Produce json
// @Success 200 {object} models.ResponsePostData
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "post data with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error
// @Failure 403 {object} models.ErrResponse "this post not belongs this creators"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Router /creators/{:creator_id}/posts/{:post_id}/{:data_id} [GET]
func (h *PostsDataIDHandler) GET(w http.ResponseWriter, r *http.Request) {
	var dataId int64
	var ok bool
	if dataId, ok = h.GetInt64FromParam(w, r, "data_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	post, err := h.postsDataUsecase.GetData(dataId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondPost := models.ToResponsePostData(*post)

	h.Log(r).Debugf("get post data with id %d", dataId)
	h.Respond(w, r, http.StatusOK, respondPost)
}

// DELETE posts_data
// @Summary delete current posts_data
// @Description delete current posts_data from current creator
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error
// @Failure 403 {object} models.ErrResponse "this post not belongs this creators"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Router /creators/{:creator_id}/posts/{:post_id}/{:data_id} [DELETE]
func (h *PostsDataIDHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	var dataId int64
	var ok bool
	if dataId, ok = h.GetInt64FromParam(w, r, "data_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	if err := h.postsDataUsecase.Delete(dataId); err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	w.WriteHeader(http.StatusOK)
}
