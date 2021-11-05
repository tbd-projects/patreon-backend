package upl_cover_posts_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type PostsUpdateCoverHandler struct {
	postsUsecase usePosts.Usecase
	bh.BaseHandler
}

func NewPostsUpdateCoverHandler(log *logrus.Logger,
	ucPosts usePosts.Usecase, manager sessions.SessionsManager) *PostsUpdateCoverHandler {
	h := &PostsUpdateCoverHandler{
		BaseHandler:  *bh.NewBaseHandler(log),
		postsUsecase: ucPosts,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)

	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)
	return h
}

// PUT CoverUpdate
// @Summary set new post cover
// @Accept  image/png, image/jpeg, image/jpg
// @Param cover formData file true "cover file with ext jpeg/png"
// @Success 200 "successfully upload cover"
// @Failure 400 {object} models.ErrResponse "size of file very big", "please upload a JPEG, JPG or PNG files", "invalid form field name for load file"
// @Failure 404 {object} models.ErrResponse "post with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 401 "User are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/update/cover [PUT]
func (h *PostsUpdateCoverHandler) PUT(w http.ResponseWriter, r *http.Request) {
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

	file, filename, code, err := h.GerFilesFromRequest(w, r, bh.MAX_UPLOAD_SIZE,
		"cover", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	err = h.postsUsecase.LoadCover(file, filename, postId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}
