package middleware

import (
	"net/http"
	hf "patreon/internal/app/delivery/http/handlers/base_handler/handler_interfaces"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	useAttaches "patreon/internal/app/usecase/attaches"
	"patreon/internal/app/utilits"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AttachesMiddleware struct {
	log             utilits.LogObject
	usecaseAttaches useAttaches.Usecase
}

func NewAttachesMiddleware(log *logrus.Logger, usecaseAttaches useAttaches.Usecase) *AttachesMiddleware {
	return &AttachesMiddleware{log: utilits.NewLogObject(log), usecaseAttaches: usecaseAttaches}
}

// CheckCorrectAttachFunc Errors
//		Status 500 middleware.InternalError
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 403 middleware.IncorrectCreatorForPost
func (mw *AttachesMiddleware) CheckCorrectAttachFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := utilits.Responder{LogObject: mw.log}
		var postId, attachId int64
		var err error

		vars := mux.Vars(r)
		id, ok := vars["attach_id"]
		attachId, err = strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		id, ok = vars["post_id"]

		postId, err = strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		attach, err := mw.usecaseAttaches.GetAttach(attachId)

		if err != nil || postId != attach.PostId {
			if err != nil && !errors.Is(err, repository.NotFound) {
				mw.log.Log(r).Errorf("some error of bd attach %v", err)
				respond.Error(w, r, http.StatusInternalServerError, BDError)
				return
			}
			mw.log.Log(r).Warnf("this attach %d not belongs to this post %d", attachId, attach.PostId)
			respond.Error(w, r, http.StatusForbidden, IncorrectAttachForPost)
			return
		}

		next(w, r)
	}
}

// CheckCorrectAttach Errors
//		Status 500 middleware.InternalError
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 403 middleware.IncorrectCreatorForPost
func (mw *AttachesMiddleware) CheckCorrectAttach(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCorrectAttachFunc(handler.ServeHTTP))
}
