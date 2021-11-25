package upd_video_attach_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	useAttaches "patreon/internal/app/usecase/attaches"
	usePosts "patreon/internal/app/usecase/posts"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type AttachUploadVideoHandler struct {
	attachesUsecase useAttaches.Usecase
	bh.BaseHandler
}

func NewAttachUploadVideoHandler(
	log *logrus.Logger,
	ucAttaches useAttaches.Usecase,
	ucPosts usePosts.Usecase,
	sClient session_client.AuthCheckerClient) *AttachUploadVideoHandler {
	h := &AttachUploadVideoHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		attachesUsecase: ucAttaches,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)
	h.AddMiddleware(sessionMiddleware.Check,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.
			NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken,
		middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost,
		middleware.NewAttachesMiddleware(log, ucAttaches).CheckCorrectAttach)

	h.AddMethod(http.MethodPut, h.PUT)

	return h
}

// PUT update video to post
// @Summary update video to post
// @tags attaches
// @Accept  video/3gpp, video/mp4
// @Param video formData file true "image file with ext video/3gpp, video/mp4"
// @Success 200
// @Failure 400 {object} http_models.ErrResponse "size of file very big", "invalid form field name for load file", "please upload a some type"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid data type", "this post id not know"
// @Failure 404 {object} http_models.ErrResponse "attach with this id not found"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/{:attach_id}/update/video [PUT]
func (h *AttachUploadVideoHandler) PUT(w http.ResponseWriter, r *http.Request) {
	var attachId int64
	var ok bool

	if attachId, ok = h.GetInt64FromParam(w, r, "attach_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
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

	err = h.attachesUsecase.UpdateVideo(file, filename, attachId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}
