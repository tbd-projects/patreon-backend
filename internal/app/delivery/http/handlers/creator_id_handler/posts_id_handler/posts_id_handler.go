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
	useUser "patreon/internal/app/usecase/user"
)

type PostsIDHandler struct {
	postsUsecase usePosts.Usecase
	userUsecase  useUser.Usecase
	bh.BaseHandler
}

func NewPostsIDHandler(log *logrus.Logger,
	ucPosts usePosts.Usecase, ucUser useUser.Usecase, manager sessions.SessionsManager) *PostsIDHandler {
	h := &PostsIDHandler{
		BaseHandler:  *bh.NewBaseHandler(log),
		postsUsecase: ucPosts,
		userUsecase:  ucUser,
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
// @Success 200 {object} http_models.ResponsePostWithAttaches "posts"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 404 {object} http_models.ErrResponse "post with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "this user not have award for this post"
// @Router /creators/{:creator_id}/posts/{:post_id} [GET]
func (h *PostsIDHandler) GET(w http.ResponseWriter, r *http.Request) {
	var postId, userId, creatorId int64
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

	if creatorId, ok = h.GetInt64FromParam(w, r, "creator_id"); !ok {
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

	access, err := h.userUsecase.CheckAccessForAward(userId, post.Awards, creatorId)

	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	if !access {
		h.Log(r).Warnf("Fobidden for user %d to post %v", userId, post)
		h.Error(w, r, http.StatusForbidden, handler_errors.UserNotHaveAward)
		return
	}

	respondPost := http_models.ToResponsePostWithAttaches(*post)

	h.Log(r).Debugf("get post with id %d", postId)
	h.Respond(w, r, http.StatusOK, respondPost)
}

// DELETE Post
// @Summary delete current post
// @Description delete current post from current creator
// @Produce json
// @Success 200 "post was delete"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
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
