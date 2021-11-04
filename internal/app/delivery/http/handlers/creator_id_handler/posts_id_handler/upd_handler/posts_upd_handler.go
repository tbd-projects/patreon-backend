package posts_upd_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	models_db "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

type PostsUpdateHandler struct {
	postsUsecase usePosts.Usecase
	bh.BaseHandler
}

func NewPostsUpdateHandler(log *logrus.Logger,
	ucPosts usePosts.Usecase, manager sessions.SessionsManager) *PostsUpdateHandler {
	h := &PostsUpdateHandler{
		BaseHandler:  *bh.NewBaseHandler(log),
		postsUsecase: ucPosts,
	}
	h.AddMiddleware(sessionMid.NewSessionMiddleware(manager, log).Check,
		middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost)

	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// PUT Posts
// @Summary update current posts
// @Description update current posts from current creator
// @Param user body models.RequestPosts true "Request body for posts"
// @Produce json
// @Success 200
// @Failure 422 {object} models.ErrResponse "invalid body in request"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "post with this id not found"
// @Failure 422 {object} models.ErrResponse "empty title"
// @Failure 422 {object} models.ErrResponse "this creator id not know"
// @Failure 422 {object} models.ErrResponse "this awards id not know"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 500 {object} models.ErrResponse "server error
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Failure 403 {object} models.ErrResponse "this post not belongs this creators"
// @Failure 401 "User are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/update [PUT]
func (h *PostsUpdateHandler) PUT(w http.ResponseWriter, r *http.Request) {
	req := &models.RequestPosts{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	postId, ok := h.GetInt64FromParam(w, r, "post_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	if err = h.postsUsecase.Update(&models_db.UpdatePost{ID: postId, Title: req.Title,
		Description: req.Description, Awards: req.AwardsId}); err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}
