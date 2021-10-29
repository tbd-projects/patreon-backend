package likes_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	"patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	useLikes "patreon/internal/app/usecase/likes"
	usePosts "patreon/internal/app/usecase/posts"
)

type LikesHandler struct {
	likesUsecase useLikes.Usecase
	bh.BaseHandler
}

func NewLikesHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	ucLikes useLikes.Usecase, ucPosts usePosts.Usecase, manager sessions.SessionsManager) *LikesHandler {
	h := &LikesHandler{
		BaseHandler:   *bh.NewBaseHandler(log, router, cors),
		likesUsecase: ucLikes,
	}
	postsMiddleware := middleware.NewPostsMiddleware(log, ucPosts)
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	creatorMiddleware := middleware.NewCreatorsMiddleware(log)
	h.AddMiddleware(sessionMiddleware.Check, creatorMiddleware.CheckAllowUser, postsMiddleware.CheckCorrectPost)
	h.AddMethod(http.MethodDelete, h.DELETE)
	h.AddMethod(http.MethodPut, h.PUT)
	return h
}

// DELETE Likes
// @Summary deletes like from set post
// @Description deletes like form post id in url
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "like with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "can not get info from context"
// @Failure 409 {object} models.ErrResponse "this user not have like for this post"
// @Failure 403 {object} models.ErrResponse "this post not belongs this creators"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Router /creators/{:creator_id}/posts/{:post_id}/like [DELETE]
func (h *LikesHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	var postsId, userId int64
	var ok bool
	postsId, ok = h.GetInt64FromParam(w, r, "post_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	userId, ok = r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ContextError)
		return
	}

	err := h.likesUsecase.Delete(postsId, userId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PUT Likes
// @Summary deletes like from set post
// @Description deletes like form post id in url
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "like with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "can not get info from context"
// @Failure 409 {object} models.ErrResponse "this user already add like for this post"
// @Failure 403 {object} models.ErrResponse "this post not belongs this creators"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Router /creators/{:creator_id}/posts/{:post_id}/like [PUT]
func (h *LikesHandler) PUT(w http.ResponseWriter, r *http.Request) {
	var postsId, userId int64
	var ok bool
	postsId, ok = h.GetInt64FromParam(w, r, "post_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	userId, ok = r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ContextError)
		return
	}

	err := h.likesUsecase.Add(&models.Like{PostId: postsId, UserId: userId})
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}
