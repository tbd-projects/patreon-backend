package comments_handler

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

type CommentsHandler struct {
	commentsUsecase useComments.Usecase
	bh.BaseHandler
}

func NewCommentsHandler(log *logrus.Logger,
	ucComments useComments.Usecase,
	ucPosts usePosts.Usecase,
	sClient client.AuthCheckerClient) *CommentsHandler {
	h := &CommentsHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		commentsUsecase: ucComments,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)

	h.AddMiddleware(sessionMiddleware.AddUserId, middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost)
	h.AddMethod(http.MethodPost, h.POST, sessionMiddleware.CheckFunc, csrf_middleware.NewCsrfMiddleware(log,
		usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)
	h.AddMethod(http.MethodGet, h.GET)
	return h
}

// POST comments
// @Summary crete current comment
// @tags comments
// @Description create current comment for current post
// @Produce json
// @Param attaches body http_models.RequestComment true "Request body for set comment"
// @Success 200 {object} http_models.IdResponse
// @Failure 400 {object} http_models.ErrResponse ""invalid parameters""
// @Failure 409 {object} http_models.ErrResponse ""comment already exist""
// @Failure 500 {object} http_models.ErrResponse ""can not do bd operation", "server error""
// @Failure 403 {object} http_models.ErrResponse ""csrf token is invalid, get new token", "this post not belongs this creators", "this user can not add comment as creator""
// @Failure 422 {object} http_models.ErrResponse ""this post id not know", "this user id not know""
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/comments [POST]
func (h *CommentsHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestComment{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	var postId, creatorId int64
	var ok bool
	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}

	if creatorId, ok = h.GetInt64FromParam(w, r, "creator_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
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

	res, err := h.commentsUsecase.Create(&models.Comment{PostId: postId, AuthorId: userID, Body: req.Body,
		AsCreator: req.AsCreator})
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	h.Log(r).Debugf("add comment to post %d", postId)
	h.Respond(w, r, http.StatusOK, http_models.IdResponse{ID: res})
}

// GET comments
// @Summary get current comments
// @tags comments
// @Description get current comments for current post
// @Param page query uint64 true "start page number of posts mutually exclusive with offset"
// @Param offset query uint64 true "start number of posts mutually exclusive with page"
// @Param limit query uint64 true "posts to return"
// @Produce json
// @Success 200 {object} http_models.ResponsePostComments
// @Failure 400 {object} http_models.ErrResponse ""invalid parameters", "invalid parameters in query""
// @Failure 409 {object} http_models.ErrResponse ""comment already exist""
// @Failure 500 {object} http_models.ErrResponse ""can not do bd operation", "server error""
// @Failure 403 {object} http_models.ErrResponse ""csrf token is invalid, get new token", "this post not belongs this creators", "this user can not add comment as creator""
// @Router /creators/{:creator_id}/posts/{:post_id}/comments [GET]
func (h *CommentsHandler) GET(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := h.GetPaginationFromQuery(w, r)
	if !ok {
		return
	}
	var postId int64
	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	res, err := h.commentsUsecase.GetPostComments(postId, &models.Pagination{Limit: limit, Offset: offset})
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	h.Log(r).Debugf("get comments from post %d", postId)
	h.Respond(w, r, http.StatusOK, http_models.ToResponsePostComments(res))
}
