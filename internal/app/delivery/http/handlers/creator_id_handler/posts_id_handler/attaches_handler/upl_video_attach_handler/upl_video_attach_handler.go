package upl_video_attach_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	http_models "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	useAttaches "patreon/internal/app/usecase/attaches"
	usePosts "patreon/internal/app/usecase/posts"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type PostsUploadVideoHandler struct {
	attachesUsecase useAttaches.Usecase
	bh.BaseHandler
}

func NewPostsUploadVideoHandler(log *logrus.Logger,
	ucAttaches useAttaches.Usecase,
	ucPosts usePosts.Usecase,
	sClient session_client.AuthCheckerClient) *PostsUploadVideoHandler {
	h := &PostsUploadVideoHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		attachesUsecase: ucAttaches,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)

	h.AddMethod(http.MethodPost, h.POST,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// POST add video to post
// @Summary add video to post
// @tags attaches
// @Accept video/3gpp, video/mp4
// @Param video formData file true "image file with ext video/3gpp, video/mp4"
// @Success 201 {object} http_models.IdResponse "id attaches"
// @Failure 400 {object} http_models.ErrResponse "size of file very big", "please upload some types", "invalid form field name for load file"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid data type", "this post id not know"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/attaches/video [POST]
func (h *PostsUploadVideoHandler) POST(w http.ResponseWriter, r *http.Request) {
	var attachId, postId int64
	var ok bool

	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}
	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	file, filename, code, err := h.GerFilesFromRequest(w, r, bh.MAX_UPLOAD_SIZE,
		"video", []string{"video/3gpp", "video/mp4"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	attachId, err = h.attachesUsecase.LoadVideo(file, filename, postId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Respond(w, r, http.StatusCreated, &http_models.IdResponse{ID: attachId})
}
