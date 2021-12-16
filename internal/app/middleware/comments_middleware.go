package middleware

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	hf "patreon/internal/app/delivery/http/handlers/base_handler/handler_interfaces"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usecase_comments "patreon/internal/app/usecase/comments"
	"patreon/internal/app/utilits"
	"strconv"
)

type CommentsMiddleware struct {
	log             utilits.LogObject
	usecaseComments usecase_comments.Usecase
}

func NewCommentsMiddleware(log *logrus.Logger, usecaseComments usecase_comments.Usecase) *CommentsMiddleware {
	return &CommentsMiddleware{log: utilits.NewLogObject(log), usecaseComments: usecaseComments}
}

// CheckCorrectCommentFunc Errors
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 404 handler_errors.CommentNotFound
//		Status 403 middleware.IncorrectCommentForPost
//		Status 403 middleware.IncorrectCommentForUser
func (mw *CommentsMiddleware) CheckCorrectCommentFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := utilits.Responder{LogObject: mw.log}
		var commentId, postId int64
		var err error

		vars := mux.Vars(r)
		id, ok := vars["post_id"]
		postId, err = strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		id, ok = vars["comment_id"]
		commentId, err = strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		comment, err := mw.usecaseComments.Get(commentId)
		if err != nil {
			if err != nil && errors.Is(err, repository.NotFound) {
				mw.log.Log(r).Warnf("comment with id %d not found", commentId)
				respond.Error(w, r, http.StatusNotFound, handler_errors.CommentNotFound)
				return
			}
			mw.log.Log(r).Errorf("some error of bd comments %v", err)
			respond.Error(w, r, http.StatusInternalServerError, BDError)
			return
		}

		if comment.PostId != postId {
			mw.log.Log(r).Warnf("comment with id %d not belongs to post with id %d", commentId, postId)
			respond.Error(w, r, http.StatusForbidden, IncorrectCommentForPost)
			return
		}

		userID, ok := r.Context().Value("user_id").(int64)
		if !ok {
			mw.log.Log(r).Error("can not get user_id from context in comment middleware")
			respond.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
			return
		}

		if comment.AuthorId != userID {
			mw.log.Log(r).Warnf("comment with id %d not belongs to this user with id %d", commentId, userID)
			respond.Error(w, r, http.StatusForbidden, IncorrectCommentForUser)
			return
		}

		next(w, r)
	}
}

// CheckCorrectComment Errors
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 404 handler_errors.CommentNotFound
//		Status 403 middleware.IncorrectCommentForPost
//		Status 403 middleware.IncorrectCommentForUser
func (mw *CommentsMiddleware) CheckCorrectComment(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCorrectCommentFunc(handler.ServeHTTP))
}
