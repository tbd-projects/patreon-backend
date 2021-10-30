package posts_id_handler

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
)

type PostsIDHandler struct {
	postsUsecase usePosts.Usecase
	bh.BaseHandler
}

func NewPostsIDHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	ucPosts usePosts.Usecase, manager sessions.SessionsManager) *PostsIDHandler {
	h := &PostsIDHandler{
		BaseHandler:  *bh.NewBaseHandler(log, router, cors),
		postsUsecase: ucPosts,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	h.AddMiddleware(middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)
	h.AddMethod(http.MethodGet, h.GET)
	return h
}

// GET Awards
// @Summary delete current awards
// @Description delete current awards from current creator
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "post with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error
// @Failure 403 {object} models.ErrResponse "this post not belongs this creators"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Router /creators/{:creator_id}/posts/{:post_id} [GET]
func (h *PostsIDHandler) GET(w http.ResponseWriter, r *http.Request) {
	var postId, userId int64
	var ok bool

	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	if userId, ok = r.Context().Value("user_id").(int64); !ok {
		userId = usePosts.EmptyUser
	}

	post, err := h.postsUsecase.GetPost(postId, userId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondPost := models.ToResponsePostWithData(*post)

	h.Log(r).Debugf("get post with id %d", postId)
	h.Respond(w, r, http.StatusOK, respondPost)
}
