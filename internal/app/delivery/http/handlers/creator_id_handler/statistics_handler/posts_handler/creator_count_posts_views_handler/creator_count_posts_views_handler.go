package statistics_count_posts_views_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	statistics_usecase "patreon/internal/app/usecase/statistics"
)

type CreatorViewsHandler struct {
	statisticsUsecase statistics_usecase.Usecase
	bh.BaseHandler
}

func NewCreatorViewsHandler(log *logrus.Logger, ucStatistics statistics_usecase.Usecase) *CreatorViewsHandler {
	h := &CreatorViewsHandler{
		statisticsUsecase: ucStatistics,
		BaseHandler:       *bh.NewBaseHandler(log),
	}

	h.AddMethod(http.MethodGet, h.GET)

	return h
}

// GET CountPostsViews
// @Summary get count of creator posts views
// @tags posts
// @Description get count of creator posts views
// @Produce json
// @Param days query uint64 true "number of processing days"
// @Success 201 {object} http_models.ResponseCreatorPostsViews
// @Failure 400 {object} http_models.ErrResponse "invalid parameters", "invalid parameters in query"
// @Failure 404 {object} http_models.ErrResponse "creator not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Router /creators/{:creator_id}/statistics/posts/views [GET]
func (h *CreatorViewsHandler) GET(w http.ResponseWriter, r *http.Request) {
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

	countPostsViews, err := h.statisticsUsecase.GetCountCreatorViews(creatorId, days)
	if err != nil {
		h.Log(r).Debugf("error happened : %v", err)
		h.UsecaseError(w, r, err, codeByErrorGet)
		return
	}

	h.Log(r).Debugf("count creator posts views = %v; creator_id = %v for last %v days", countPostsViews, creatorId, days)
	h.Respond(w, r, http.StatusOK, http_models.ResponseCreatorPostsViews{CountPostsViews: countPostsViews})
}
