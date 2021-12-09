package comments_id_handler

import (
	"fmt"
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	"patreon/internal/app/models"
	useComments "patreon/internal/app/usecase/comments"
	usePosts "patreon/internal/app/usecase/posts"
	"patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/microcosm-cc/bluemonday"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CommentsIdHandler struct {
	commentsUsecase useComments.Usecase
	bh.BaseHandler
}

func NewCommentsIdHandler(log *logrus.Logger,
	ucComments useComments.Usecase,
	ucPosts usePosts.Usecase,
	sClient client.AuthCheckerClient) *CommentsIdHandler {
	h := &CommentsIdHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		commentsUsecase: ucComments,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)

	h.AddMiddleware(sessionMiddleware.Check, csrf_middleware.NewCsrfMiddleware(log,
		usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost,
		middleware.NewCommentsMiddleware(log, ucComments).CheckCorrectComment)
	h.AddMethod(http.MethodPut, h.PUT)
	h.AddMethod(http.MethodDelete, h.DELETE)
	return h
}

// PUT comments
// @Summary update current comment
// @tags comments
// @Description update current comment for current post
// @Produce json
// @Param attaches body http_models.RequestComment true "Request body for update comment"
// @Success 200
// @Failure 400 {object} http_models.ErrResponse ""invalid parameters""
// @Failure 404 {object} http_models.ErrResponse ""comment with this id not found""
// @Failure 500 {object} http_models.ErrResponse ""can not do bd operation", "server error""
// @Failure 403 {object} http_models.ErrResponse ""this comment not belongs this post", "this comment not belongs this user", "csrf token is invalid, get new token", "this user can not add comment as creator", "this post not belongs this creators""
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/comments/{:comment_id} [PUT]
func (h *CommentsIdHandler) PUT(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestComment{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	var commentId, creatorId int64
	var ok bool
	if commentId, ok = h.GetInt64FromParam(w, r, "comment_id"); !ok {
		return
	}

	if creatorId, ok = h.GetInt64FromParam(w, r, "creator_id"); !ok {
		return
	}


	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	if req.AsCreator && userID != creatorId {
		h.Log(r).Error(fmt.Sprintf("try add as creator for not self post: creatorId %d, userId %d",
			userID, creatorId))
		h.Error(w, r, http.StatusForbidden, handler_errors.NotAllowAddComment)
		return
	}


	err = h.commentsUsecase.Update(&models.Comment{ID: commentId, Body: req.Body, AsCreator: req.AsCreator})
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPUT)
		return
	}

	h.Log(r).Debugf("update comment with id %d", commentId)
	w.WriteHeader(http.StatusOK)
}

// DELETE comments
// @Summary update current comment
// @tags comments
// @Description update current comment from current post
// @Produce json
// @Success 200
// @Failure 400 {object} http_models.ErrResponse ""invalid parameters""
// @Failure 404 {object} http_models.ErrResponse ""comment with this id not found""
// @Failure 500 {object} http_models.ErrResponse ""can not do bd operation", "server error""
// @Failure 403 {object} http_models.ErrResponse ""this comment not belongs this post", "this comment not belongs this user", "csrf token is invalid, get new token", "this post not belongs this creators""
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/comments/{:comment_id} [DELETE]
func (h *CommentsIdHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	var commentId int64
	var ok bool
	if commentId, ok = h.GetInt64FromParam(w, r, "comment_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	err := h.commentsUsecase.Delete(commentId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}

	h.Log(r).Debugf("delete comment with id %d", commentId)
	w.WriteHeader(http.StatusOK)
}
