package middleware

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	hf "patreon/internal/app/delivery/http/handlers/base_handler/handler_interfaces"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usecase_awards "patreon/internal/app/usecase/awards"
	"patreon/internal/app/utilits"
	"strconv"
)

type AwardsMiddleware struct {
	log           utilits.LogObject
	usecaseAwards usecase_awards.Usecase
}

func NewAwardsMiddleware(log *logrus.Logger, usecaseAwards usecase_awards.Usecase) *AwardsMiddleware {
	return &AwardsMiddleware{log: utilits.NewLogObject(log), usecaseAwards: usecaseAwards}
}

// CheckCorrectAwardFunc Errors
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 403 middleware.IncorrectCreatorForAward
func (mw *AwardsMiddleware) CheckCorrectAwardFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := utilits.Responder{LogObject: mw.log}
		var awardsId, cretorId, bdCretorId int64
		var err error

		vars := mux.Vars(r)
		id, ok := vars["creator_id"]
		cretorId, err = strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		id, ok = vars["award_id"]

		awardsId, err = strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		bdCretorId, err = mw.usecaseAwards.GetCreatorId(awardsId)

		if err != nil || bdCretorId != cretorId {
			if err != nil && !errors.Is(err, repository.NotFound) {
				mw.log.Log(r).Errorf("some error of bd awards %v", err)
				respond.Error(w, r, http.StatusInternalServerError, BDError)
				return
			}
			mw.log.Log(r).Warnf("this post %d not belongs to this creator %d", awardsId, cretorId)
			respond.Error(w, r, http.StatusForbidden, IncorrectCreatorForAward)
			return
		}

		next(w, r)
	}
}

func (mw *AwardsMiddleware) CheckCorrectAward(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCorrectAwardFunc(handler.ServeHTTP))
}
