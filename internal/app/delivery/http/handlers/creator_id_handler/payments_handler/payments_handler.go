package creator_payments_handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	db_models "patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase/payments"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"
	"strconv"

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
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
	)
	return h
}

func (h *PaymentsHandler) redirect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("page", "1")
	query.Set("limit", fmt.Sprintf("%d", payments.BaseLimit))
	r.URL.RawQuery = query.Encode()
	redirectUrl := r.URL.String()
	h.Log(r).Debugf("redirect to url: %s, with offest 0 and limit %d", redirectUrl, payments.BaseLimit)

	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

// GET CreatorPayments
// @Summary get all creator payments
// @tags payments
// @Description get all creator payments
// @Produce json
// @Param page query uint64 true "start page number of posts mutually exclusive with offset"
// @Param offset query uint64 true "start number of posts mutually exclusive with page"
// @Param limit query uint64 true "posts to return"
// @Success 200 {object} http_models.ResponseCreatorPayments "Success"
// @Failure 204 {object} http_models.OkResponse "payments not found"
// @Failure 500 {object} http_models.ErrResponse "server error"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/payments [GET]
func (h *PaymentsHandler) GET(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := h.GetPaginationFromQuery(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["creator_id"]
	creatorID, err := strconv.ParseInt(id, 10, 64)
	if !ok || err != nil {
		h.Log(r).Infof("invalid parametrs creator_id %v", vars)
		h.Respond(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	creatorPayments, err := h.paymentsUsecase.GetCreatorPayments(creatorID,
		&db_models.Pagination{Limit: limit, Offset: offset})
	if err != nil {
		if err == repository.NotFound {
			h.Respond(w, r, http.StatusNoContent, http_models.OkResponse{
				Ok: handler_errors.PaymentsNotFound.Error(),
			})
		} else {
			h.UsecaseError(w, r, err, codeByErrorGET)
		}
		return
	}
	res := http_models.ToResponseCreatorPayments(creatorPayments)

	h.Respond(w, r, http.StatusOK, res)
}
