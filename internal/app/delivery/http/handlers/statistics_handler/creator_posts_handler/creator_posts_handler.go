package statistics_handler

import (
	"github.com/sirupsen/logrus"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	statistics_usecase "patreon/internal/app/usecase/statistics"
)

type CreatorPostsHandler struct {
	statisticsUsecase statistics_usecase.Usecase
	bh.BaseHandler
}

func NewCreatorPostsHandler(log *logrus.Logger, ucStatistics statistics_usecase.Usecase) *CreatorPostsHandler {
	h := &CreatorPostsHandler{
		statisticsUsecase: ucStatistics,
		BaseHandler:       *bh.NewBaseHandler(log),
	}
	return h
}
