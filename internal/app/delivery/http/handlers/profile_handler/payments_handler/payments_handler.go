package payments_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase/payments"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/sirupsen/logrus"
)

type PaymentsHandler struct {
	sessionClient   session_client.AuthCheckerClient
	paymentsUsecase payments.Usecase
	bh.BaseHandler
}

func NewPaymentsHandler(log *logrus.Logger,
	sClient session_client.AuthCheckerClient, ucPayments payments.Usecase) *PaymentsHandler {
	h := &PaymentsHandler{
		sessionClient:   sClient,
		paymentsUsecase: ucPayments,
		BaseHandler:     *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)
	return h
}

// GET UserPayments
// @Summary get all user payments
// @Description get all user payments
// @Produce json
// @Success 200 {object} models.ResponseUserPayments "Success"
// @Failure 204 {object} models.OkResponse "payments not Found"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 401 "user are not authorized"
// @Router /user/payments [GET]
func (h *PaymentsHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	userPayments, err := h.paymentsUsecase.GetUserPayments(userID.(int64))
	if err != nil {
		if err == repository.NotFound {
			h.Respond(w, r, http.StatusNoContent, models.OkResponse{
				Ok: handler_errors.PaymentsNotFound.Error(),
			})
		} else {
			h.UsecaseError(w, r, err, codeByErrorGET)
		}
		return
	}
	res := models.ToResponseUserPayments(userPayments)

	h.Respond(w, r, http.StatusOK, res)
}
