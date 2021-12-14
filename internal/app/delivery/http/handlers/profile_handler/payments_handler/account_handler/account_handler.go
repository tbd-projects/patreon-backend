package payments_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/models"
	usecase_pay_token "patreon/internal/app/usecase/pay_token"
)

type AccountHandler struct {
	tokenUsecase usecase_pay_token.Usecase
	bh.BaseHandler
}

func NewAccountHandler(log *logrus.Logger, ucPayToken usecase_pay_token.Usecase) *AccountHandler {
	h := &AccountHandler{
		tokenUsecase: ucPayToken,
		BaseHandler:  *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET)

	return h
}

// GET Account number
// @Summary get number of yoomoney account
// @tags payments
// @Description get number of yoomoney account
// @Produce json
// @Success 200 {object} http_models.ResponsePayAccount "Success"
// @Failure 500 {object} http_models.ErrResponse "server error"
// @Router /user/payments/account [GET]
func (h *AccountHandler) GET(w http.ResponseWriter, r *http.Request) {
	account := h.tokenUsecase.GetAccount()

	h.Respond(w, r, http.StatusOK, models.PayAccount{Account: account})
}
