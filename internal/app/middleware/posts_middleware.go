package middleware

import (
	"net/http"
	hf "patreon/internal/app/delivery/http/handlers/base_handler/handler_interfaces"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usePosts "patreon/internal/app/usecase/posts"
	"patreon/internal/app/utilits"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PostsMiddleware struct {
	log          utilits.LogObject
	usecasePosts usePosts.Usecase
}

func NewPostsMiddleware(log *logrus.Logger, usecasePosts usePosts.Usecase) *PostsMiddleware {
	return &PostsMiddleware{log: utilits.NewLogObject(log), usecasePosts: usecasePosts}
}

// CheckCorrectPostFunc Errors
//		Status 500 middleware.InternalError
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 403 middleware.IncorrectCreatorForPost
func (mw *PostsMiddleware) CheckCorrectPostFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := utilits.Responder{LogObject: mw.log}
		var postId, creatorId, bdCreatorId int64
		var err error

		vars := mux.Vars(r)
		id, ok := vars["creator_id"]
		creatorId, err = strconv.ParseInt(id, 10, 64)
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

		bdCreatorId, err = mw.usecasePosts.GetCreatorId(postId)

		if err != nil || bdCreatorId != creatorId {
			if err != nil && !errors.Is(err, repository.NotFound) {
				mw.log.Log(r).Errorf("some error of bd awards %v", err)
				respond.Error(w, r, http.StatusInternalServerError, BDError)
				return
			}
			mw.log.Log(r).Warnf("this post %d not belongs to this creator %d", postId, creatorId)
			respond.Error(w, r, http.StatusForbidden, IncorrectCreatorForAward)
			return
		}

		next(w, r)
	}
}

// CheckCorrectPost Errors
//		Status 500 middleware.InternalError
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 403 middleware.IncorrectCreatorForPost
func (mw *PostsMiddleware) CheckCorrectPost(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCorrectPostFunc(handler.ServeHTTP))
}
