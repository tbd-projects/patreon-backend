package payments_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	usecase_pay_token "patreon/internal/app/usecase/pay_token"
	"patreon/internal/app/usecase/payments"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/sirupsen/logrus"
)

type TokenHandler struct {
	sessionClient   session_client.AuthCheckerClient
	tokenUsecase    usecase_pay_token.Usecase
	paymentsUsecase payments.Usecase
	bh.BaseHandler
}

func NewTokenHandler(log *logrus.Logger,
	sClient session_client.AuthCheckerClient, ucPayToken usecase_pay_token.Usecase, ucPayments payments.Usecase) *TokenHandler {
	h := &TokenHandler{
		sessionClient:   sClient,
		tokenUsecase:    ucPayToken,
		paymentsUsecase: ucPayments,
		BaseHandler:     *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)
	h.AddMethod(http.MethodGet, h.POST,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// GET PayToken
// @Summary get token for payments
// @tags payments
// @Description get token for payments
// @Produce json
// @Success 200 {object} http_models.ResponsePayToken "Success"
// @Failure 500 {object} http_models.ErrResponse "server error"
// @Failure 401 "user are not authorized"
// @Router /user/payments/token [GET]
func (h *TokenHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	payToken, err := h.tokenUsecase.GetToken(userID.(int64))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Respond(w, r, http.StatusOK, models.PayToken{Token: payToken.Token})
}
func (h *TokenHandler) POST(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	if err := r.ParseForm(); err != nil {
		h.Log(r).Error("can not parse form")
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidBody)
		return
	}
	payToken := r.PostForm["label"][0]

	exists, err := h.tokenUsecase.CheckToken(models.PayToken{Token: payToken})
	if !exists {
		h.UsecaseError(w, r, err, codeByErrorPOST)
		return
	}
	if err = h.paymentsUsecase.UpdateStatus(payToken); err != nil {
		h.UsecaseError(w, r, err, codeByErrorPOST)
		return
	}
	w.WriteHeader(http.StatusOK)
}
