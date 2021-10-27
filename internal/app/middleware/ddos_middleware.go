package middleware

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	usecase_access "patreon/internal/app/usecase/access"

	"github.com/sirupsen/logrus"
)

type DDosMiddleware struct {
	log           *logrus.Logger
	accessUsecase usecase_access.Usecase
}

func NewDdosMiddleware(log *logrus.Logger, accessUc usecase_access.Usecase) DDosMiddleware {
	return DDosMiddleware{
		log:           log,
		accessUsecase: accessUc,
	}
}
func (mw *DDosMiddleware) CheckAccess(next bh.HandlerFunc) bh.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIp := r.RemoteAddr
		ok, err := mw.accessUsecase.CheckBlackList(userIp)
		if ok {
			mw.log.Infof("DDOS_Middleware user with ip: %v in blackList", userIp)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		ok, err = mw.accessUsecase.CheckAccess(userIp)
		if !ok {
			if err == usecase_access.NoAccess {
				_ = mw.accessUsecase.AddToBlackList(userIp)
				mw.log.Infof("DDOS_Middleware user with ip: %v add in blackList", userIp)
				w.WriteHeader(http.StatusTooManyRequests)
				return
			} else {
				mw.log.Infof("DDOS_Middleware error user with ip: %v err: %v",
					userIp, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if err == usecase_access.FirstQuery {
			ok, err = mw.accessUsecase.Create(userIp)
			if err != nil {
				mw.log.Infof("DDOS_Middleware - error on create AccessUserCounter from user with ip: %v err: %v",
					userIp, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res, err := mw.accessUsecase.Update(userIp)
			if err != nil {
				mw.log.Infof("DDOS_Middleware - error on add user with ip: %v to BlackList err: %v",
					userIp, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			mw.log.Infof("DDOS_Middleware - Count of querys: %v from userIp: %v", res, userIp)
			next(w, r)
		}
	}
}
