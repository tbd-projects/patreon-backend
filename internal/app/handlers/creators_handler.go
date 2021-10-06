package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/handlers/base_handler"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/models"

	"github.com/sirupsen/logrus"
)

type CreatorHandler struct {
	dataStorage app.DataStorage
	bh.BaseHandler
}

func NewCreatorHandler(log *logrus.Logger, storage app.DataStorage) *CreatorHandler {
	h := &CreatorHandler{
		BaseHandler: *bh.NewBaseHandler(log),
		dataStorage: storage,
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.dataStorage.SessionManager(), h.Log()).Check)
	return h
}

// Creators
// @Summary get list of Creators
// @Description get list of creators which register on service
// @Produce json
// @Success 201 {array} models.ResponseCreator "Create user successfully"
// @Failure 503 {object} models.BaseResponse "Internal error"
// @Failure 418 "User are not authorized"
// @Router /creators [GET]
func (h *CreatorHandler) GET(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log().Error(err)
		}
	}(r.Body)
	creators, err := h.dataStorage.Store().Creator().GetCreators()
	if err != nil {
		h.Log().Errorf("get: %v err:%v can not get user from db", creators, err)
		h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}
	respondCreators := make([]models.ResponseCreator, len(creators))
	for i, cr := range creators {
		respondCreators[i] = models.ToResponseCreator(cr)
	}

	h.Log().Debugf("get creators %v", respondCreators)
	h.Respond(w, r, http.StatusOK, respondCreators)
}
