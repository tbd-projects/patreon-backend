package payments_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	usecase_pay_token "patreon/internal/app/usecase/pay_token"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/sirupsen/logrus"
)

type TokenHandler struct {
	sessionClient session_client.AuthCheckerClient
	tokenUsecase  usecase_pay_token.Usecase
	bh.BaseHandler
}

func NewTokenHandler(log *logrus.Logger,
	sClient session_client.AuthCheckerClient, ucPayToken usecase_pay_token.Usecase) *TokenHandler {
	h := &TokenHandler{
		sessionClient: sClient,
		tokenUsecase:  ucPayToken,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
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
