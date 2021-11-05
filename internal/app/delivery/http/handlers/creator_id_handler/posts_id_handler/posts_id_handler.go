package posts_id_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"
)

type PostsIDHandler struct {
	postsUsecase usePosts.Usecase
	bh.BaseHandler
}

func NewPostsIDHandler(log *logrus.Logger,
	ucPosts usePosts.Usecase, manager sessions.SessionsManager) *PostsIDHandler {
	h := &PostsIDHandler{
		BaseHandler:  *bh.NewBaseHandler(log),
		postsUsecase: ucPosts,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	postMid := middleware.NewPostsMiddleware(log, ucPosts)
	h.AddMethod(http.MethodGet, h.GET, postMid.CheckCorrectPostFunc, sessionMiddleware.AddUserIdFunc)
	h.AddMethod(http.MethodDelete, h.DELETE, sessionMiddleware.CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc, postMid.CheckCorrectPostFunc)
	return h
}

// GET Post
// @Summary get current post
// @Description get current post from current creator
// @Produce json
// @Param add-view query string false "IMPORTANT: value yes or no, - if need add view to this post"
// @Success 200 {object} models.ResponsePostWithData "posts"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "post with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators"
// @Router /creators/{:creator_id}/posts/{:post_id} [GET]
func (h *PostsIDHandler) GET(w http.ResponseWriter, r *http.Request) {
	var postId, userId int64
	var addView bool
	var ok bool

	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	value := r.URL.Query().Get("add-view")
	if value == "" {
		addView = false
	} else {
		addView = value == "yes"
	}

	if userId, ok = r.Context().Value("user_id").(int64); !ok {
		userId = usePosts.EmptyUser
	}

	post, err := h.postsUsecase.GetPost(postId, userId, addView)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondPost := models.ToResponsePostWithData(*post)

	h.Log(r).Debugf("get post with id %d", postId)
	h.Respond(w, r, http.StatusOK, respondPost)
}

// DELETE Post
// @Summary delete current post
// @Description delete current post from current creator
// @Produce json
// @Success 200 "post was delete"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id} [DELETE]
func (h *PostsIDHandler) DELETE(w http.ResponseWriter, r *http.Request) {
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

	err := h.postsUsecase.Delete(postId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}

	h.Log(r).Debugf("delete post with id %d", postId)
	w.WriteHeader(http.StatusOK)
}