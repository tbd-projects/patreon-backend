package upl_img_data_handler

import (
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
	usePostsData "patreon/internal/app/usecase/posts_data"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type PostsUploadImageHandler struct {
	postsDataUsecase usePostsData.Usecase
	bh.BaseHandler
}

func NewPostsUploadImageHandler(log *logrus.Logger,
	ucPostsData usePostsData.Usecase, ucPosts usePosts.Usecase, manager sessions.SessionsManager) *PostsUploadImageHandler {
	h := &PostsUploadImageHandler{
		BaseHandler:      *bh.NewBaseHandler(log),
		postsDataUsecase: ucPostsData,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)

	h.AddMethod(http.MethodPost, h.POST,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// POST add image to post
// @Summary add image to post
// @Accept  image/png, image/jpeg, image/jpg
// @Param image formData file true "image file with ext jpeg/png"
// @Success 201 {object} models.IdResponse "id posts_data"
// @Failure 400 {object} models.ErrResponse "size of file very big", "please upload a JPEG, JPG or PNG files", "invalid form field name for load file"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} models.ErrResponse "invalid data type", "this post id not know"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "User are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/image [POST]
func (h *PostsUploadImageHandler) POST(w http.ResponseWriter, r *http.Request) {
	var dataId, postId int64
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
		"image", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	dataId, err = h.postsDataUsecase.LoadImage(file, filename, postId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Respond(w, r, http.StatusCreated, &models.IdResponse{ID: dataId})
}
