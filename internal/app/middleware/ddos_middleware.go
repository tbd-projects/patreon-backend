package middleware

import (
	"net/http"
	usecase_access "patreon/internal/app/usecase/access"
	"patreon/internal/app/utilits"

	"github.com/sirupsen/logrus"
)

type DDosMiddleware struct {
	utilits.LogObject
	accessUsecase usecase_access.Usecase
}

func NewDdosMiddleware(log *logrus.Logger, accessUc usecase_access.Usecase) DDosMiddleware {
	return DDosMiddleware{
		LogObject:     utilits.NewLogObject(log),
		accessUsecase: accessUc,
	}
}
func (mw *DDosMiddleware) CheckAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUrl := r.URL.Path
		userIp := r.RemoteAddr + userUrl
		ok, err := mw.accessUsecase.CheckBlackList(userIp)
		if ok || err != nil {
			mw.Log(r).Warnf("DDOS_Middleware user with ip: %v in blackList", userIp)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		ok, err = mw.accessUsecase.CheckAccess(userIp)
		if !ok {
			if err == usecase_access.NoAccess {
				_ = mw.accessUsecase.AddToBlackList(userIp)
				mw.Log(r).Infof("DDOS_Middleware user with ip: %v add in blackList", userIp)
				w.WriteHeader(http.StatusTooManyRequests)
				return
			} else {
				mw.Log(r).Errorf("DDOS_Middleware error user with ip: %v err: %v",
					userIp, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if err == usecase_access.FirstQuery {
			ok, err = mw.accessUsecase.Create(userIp)
			if err != nil || !ok {
				mw.Log(r).Errorf("DDOS_Middleware - error on create AccessUserCounter from user with ip: %v err: %v",
					userIp, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res, err := mw.accessUsecase.Update(userIp)
			if err == usecase_access.NoAccess {
				_ = mw.accessUsecase.AddToBlackList(userIp)
				mw.Log(r).Infof("DDOS_Middleware user with ip: %v add in blackList", userIp)
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			if err != nil {
				mw.Log(r).Errorf("DDOS_Middleware - error on add user with ip: %v to BlackList err: %v",
					userIp, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			mw.Log(r).Infof("DDOS_Middleware - Count of querys: %v from userIp: %v", res, userIp)
		}
		next.ServeHTTP(w, r)
	})
}
