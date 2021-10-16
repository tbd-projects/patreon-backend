package creator_handler

import (
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_creator "patreon/internal/app/usecase/creator"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type CreatorHandler struct {
	sessionManager sessions.SessionsManager
	creatorUsecase usecase_creator.Usecase
	bh.BaseHandler
}

func NewCreatorHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucCreator usecase_creator.Usecase) *CreatorHandler {
	h := &CreatorHandler{
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
		creatorUsecase: ucCreator,
		sessionManager: sManager,
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check)
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
			h.Log(r).Error(err)
		}
	}(r.Body)
	creators, err := h.creatorUsecase.GetCreators()
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrors)
		return
	}

	respondCreators := make([]models.ResponseCreator, len(creators))
	for i, cr := range creators {
		respondCreators[i] = models.ToResponseCreator(cr)
	}

	h.Log(r).Debugf("get creators %v", respondCreators)
	h.Respond(w, r, http.StatusOK, respondCreators)
}
