package middleware

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	hf "patreon/internal/app/delivery/http/handlers/base_handler/handler_interfaces"
	"patreon/internal/app/utilits"
	"strconv"
)

type CreatorsMiddleware struct {
	log utilits.LogObject
}

func NewCreatorsMiddleware(log *logrus.Logger) *CreatorsMiddleware {
	return &CreatorsMiddleware{log: utilits.NewLogObject(log)}
}

// CheckAllowUserFunc Errors
//		Status 500 middleware.ContextError
//		Status 400 middleware.InvalidParameters
//		Status 403 middleware.ForbiddenChangeCreator
func (mw *CreatorsMiddleware) CheckAllowUserFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := utilits.Responder{LogObject: mw.log}
		userId := r.Context().Value("user_id")
		if userId == nil {
			mw.log.Log(r).Error("can not get user_id from context")
			respond.Error(w, r, http.StatusInternalServerError, ContextError)
			return
		}

		vars := mux.Vars(r)
		id, ok := vars["creator_id"]
		idInt, err := strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Infof("invalid parametrs creator_id %d", vars)
			respond.Error(w, r, http.StatusBadRequest, InvalidParameters)
			return
		}

		if idInt != userId {
			mw.log.Log(r).Warnf("forbidden change by user %d creator %d", userId, idInt)
			respond.Error(w, r, http.StatusForbidden, ForbiddenChangeCreator)
			return
		}

		next(w, r)
	}
}

func (mw *CreatorsMiddleware) CheckAllowUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckAllowUserFunc(handler.ServeHTTP))
}
