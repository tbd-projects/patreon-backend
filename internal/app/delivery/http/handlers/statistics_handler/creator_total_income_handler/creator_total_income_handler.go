package statistics_handler

import (
	"github.com/sirupsen/logrus"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
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
