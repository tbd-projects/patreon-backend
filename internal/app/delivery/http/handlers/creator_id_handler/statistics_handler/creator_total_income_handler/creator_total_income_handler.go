package statistics_total_income_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	http_models "patreon/internal/app/delivery/http/models"
	statistics_usecase "patreon/internal/app/usecase/statistics"
)

type CreatorTotalIncomeHandler struct {
	statisticsUsecase statistics_usecase.Usecase
	bh.BaseHandler
}

func NewCreatorTotalIncomeHandler(log *logrus.Logger, ucStatistics statistics_usecase.Usecase) *CreatorTotalIncomeHandler {
	h := &CreatorTotalIncomeHandler{
		statisticsUsecase: ucStatistics,
		BaseHandler:       *bh.NewBaseHandler(log),
	}
	return h
}

// GET CreatorTotalIncome
// @Summary get creator total income
// @tags posts
// @Description get creator total income
// @Produce json
// @Param days query uint64 true "number of processing days"
// @Success 201 {array} http_models.ResponseCreatorTotalIncome
// @Failure 400 {object} http_models.ErrResponse "invalid parameters", "invalid parameters in query"
// @Failure 404 {object} http_models.ErrResponse "creator not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Router /creators/{:creator_id}/statistics/total_income [GET]
func (h *CreatorTotalIncomeHandler) GET(w http.ResponseWriter, r *http.Request) {
	days, ok := h.GetInt64FromQueries(w, r, "days")
	if !ok {
		return
	}
	if days < 0 {
		h.Log(r).Infof("query param days < 0; days from query =  %v)", days)
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	var creatorId int64

	creatorId, ok = h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	totalIncome, err := h.statisticsUsecase.GetTotalIncome(creatorId, days)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGet)
		return
	}

	h.Log(r).Debugf("number of creator total income = %v with creator_id = %v with last %v days", totalIncome, creatorId, days)
	h.Respond(w, r, http.StatusOK, http_models.ResponseCreatorTotalIncome{TotalIncome: totalIncome})
}
