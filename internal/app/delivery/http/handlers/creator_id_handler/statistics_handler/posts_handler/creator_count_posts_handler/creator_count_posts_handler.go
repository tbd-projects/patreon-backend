package statistics_count_posts_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	statistics_usecase "patreon/internal/app/usecase/statistics"
)

type CreatorCountPostsHandler struct {
	statisticsUsecase statistics_usecase.Usecase
	bh.BaseHandler
}

func NewCreatorCountPostsHandler(log *logrus.Logger, ucStatistics statistics_usecase.Usecase) *CreatorCountPostsHandler {
	h := &CreatorCountPostsHandler{
		statisticsUsecase: ucStatistics,
		BaseHandler:       *bh.NewBaseHandler(log),
	}
	return h
}

// GET CountPosts
// @Summary get count of creator posts
// @tags posts
// @Description get count of creator posts
// @Produce json
// @Success 201 {object} http_models.ResponseCreatorCountPosts
// @Failure 400 {object} http_models.ErrResponse "invalid parameters", "invalid parameters in query"
// @Failure 404 {object} http_models.ErrResponse "creator not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Router /creators/{:creator_id}/statistics/posts/count [GET]
func (h *CreatorCountPostsHandler) GET(w http.ResponseWriter, r *http.Request) {
	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	var creatorId int64

	creatorId, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	countPostsViews, err := h.statisticsUsecase.GetCountCreatorPosts(creatorId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGet)
		return
	}

	h.Log(r).Debugf("count creator posts = %v with creator_id = %v", countPostsViews, creatorId)
	h.Respond(w, r, http.StatusOK, http_models.ResponseCreatorCountPosts{CountPosts: countPostsViews})
}
