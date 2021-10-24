package avatar_handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_user "patreon/internal/app/usecase/user"
	"sort"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 4 // 4MB

type UpdateAvatarHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewUpdateAvatarHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	sManager sessions.SessionsManager, ucUser usecase_user.Usecase) *UpdateAvatarHandler {
	h := &UpdateAvatarHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
	}
	h.AddMethod(http.MethodPut, h.PUT)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check)
	return h
}

// PUT AvatarChange
// @Summary set new user avatar
// @Accept  mpfd
// @Produce json
// @Param file formData file true "Avatar file with ext jpeg/png"
// @Success 200 "successfully upload avatar"
// @Failure 400 {object} models.ErrResponse "size of file very big"
// @Failure 400 {object} models.ErrResponse "invalid form field name"
// @Failure 400 {object} models.ErrResponse "invalid avatar extension"
// @Failure 500 {object} models.ErrResponse "server error"
// @Router /user/update/avatar [PUT]
func (h *UpdateAvatarHandler) PUT(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		h.HandlerError(w, r, http.StatusBadRequest, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.FileSizeError,
		})
		return
	}
	f, fHeader, err := r.FormFile("avatar")
	if err != nil {
		h.HandlerError(w, r, http.StatusBadRequest, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InvalidFormFieldName,
		})
		return
	}
	buff := make([]byte, 512)
	if _, err = f.Read(buff); err != nil {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InternalError,
		})
		return
	}
	validFileTypes := []string{"image/png", "image/jpeg"}
	fType := http.DetectContentType(buff)
	if pos := sort.SearchStrings(validFileTypes, fType); pos == len(validFileTypes) {
		h.HandlerError(w, r, http.StatusBadRequest, handler_errors.InvalidAvatarExt)
		return
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InternalError,
		})
		return
	}
	rootPath, err := os.Getwd()
	if err != nil {
		h.HandlerError(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			Err:         handler_errors.InternalError,
			ExternalErr: errors.New("context parse userId error"),
		})
		return
	}
	avatarPath := rootPath + "internal/media/img/" + strconv.Itoa(int(userId)) + filepath.Ext(fHeader.Filename)
	err = os.MkdirAll("internal/media/img/", os.ModePerm)
	if err != nil {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			Err:         handler_errors.InternalError,
			ExternalErr: err,
		})
		return
	}
	dst, err := os.OpenFile(avatarPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			Err:         handler_errors.InternalError,
			ExternalErr: err,
		})
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(dst)

	if _, err = io.Copy(dst, f); err != nil {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			Err:         handler_errors.InternalError,
			ExternalErr: err,
		})
		return
	}
	err = h.userUsecase.UpdateAvatar(userId, avatarPath)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
	}
	w.WriteHeader(http.StatusOK)
}
