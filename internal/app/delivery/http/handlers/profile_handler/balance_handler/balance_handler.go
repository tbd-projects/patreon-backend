package balance_handler

import (
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	responseModels "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_user "patreon/internal/app/usecase/user"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type BalanceHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewBalanceHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	sManager sessions.SessionsManager, ucUser usecase_user.Usecase) *BalanceHandler {
	h := &BalanceHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
	}
	h.AddMethod(http.MethodGet, h.GET,
		middleware.NewSessionMiddleware(h.sessionManager, log).CheckFunc,
	)
	return h
}

// GET Balance
// @Summary get information about user balance
// @Description get the user balance
// @Produce json
// @Success 200 {object} models.User "Success"
// @Failure 401 "User are not authorized"
// @Failure 404 {object} models.ErrResponse "user with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Router /user/balance [GET]
func (h *BalanceHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	balance, err := h.userUsecase.GetBalance(userID.(int64))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Log(r).Debugf("get user balance %v", balance)
	h.Respond(w, r, http.StatusOK, responseModels.ResponseBalance{ID: userID.(int64), Balance: balance})
}

// PUT Balance
// @Summary top up the user balance
// @Description top up the user balance and get updated
// @Accept  models.RequestBalanceTop
// @Produce models.ResponseBalance
// @Success 200 {object} models.ResponseBalance "Success"
// @Failure 401 "User are not authorized"
// @Failure 404 {object} models.ErrResponse "user with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Router /user/balance/top-up [GET]
func (h *BalanceHandler) PUT(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	balance, err := h.userUsecase.GetBalance(userID.(int64))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Log(r).Debugf("get user balance %v", balance)
	h.Respond(w, r, http.StatusOK, responseModels.ResponseBalance{ID: userID.(int64), Balance: balance})
}
