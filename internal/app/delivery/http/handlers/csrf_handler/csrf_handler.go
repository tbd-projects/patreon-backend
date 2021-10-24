package csrf_handler

import (
	"io"
	"net/http"
	"patreon/internal/app"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/sessions"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type CsrfHandler struct {
	csrfUsecase    usecase_csrf.Usecase
	sessionManager sessions.SessionsManager
	bh.BaseHandler
}

func NewCsrfHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig, sManager sessions.SessionsManager) *CsrfHandler {
	h := &CsrfHandler{
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
		sessionManager: sManager,
	}
	h.AddMethod(http.MethodGet, h.GET)
	return h
}

// GET CSRF Token
// @Summary get CSRF Token
// @Description generate usecase token and return to client
// @Produce json
// @Success 201 {array} models.ResponseCreator
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Router /creators [GET]
func (h *CsrfHandler) GET(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	//h.Log(r).Debugf("get creators %v", respondCreators)
	//h.Respond(w, r, http.StatusOK, respondCreators)
}
