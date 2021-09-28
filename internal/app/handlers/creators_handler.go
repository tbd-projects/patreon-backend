package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"
	"patreon/internal/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CreatorHandler struct {
	baseHandler    app.HandlerJoiner
	authMiddleware middleware.SessionMiddleware
	Store          store.Store
	SessionManager sessions.SessionsManager
	RespondHandler
}

func NewCreatorHandler() *CreatorHandler {
	return &CreatorHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/creators"),
		RespondHandler: RespondHandler{logrus.New()},
	}
}

func (h *CreatorHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *CreatorHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.log)
}
func (h *CreatorHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.Check(h)).Methods("GET", "OPTIONS")
	h.baseHandler.Join(router)
}
func (h *CreatorHandler) JoinHandlers(joinedHandlers []app.Joinable) {
	h.baseHandler.AddHandlers(joinedHandlers)
}

// Creators
// @Summary get list of Creators
// @Description get list of creators which register on service
// @Produce json
// @Success 201 {array} models.ResponseCreator "Create user successfully"
// @Failure 503 {object} models.BaseResponse "Internal error"
// @Failure 418 "User are not authorized"
// @Router /creators [GET]
func (h *CreatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Error(err)
		}
	}(r.Body)
	creators, err := h.Store.Creator().GetCreators()
	if err != nil {
		h.log.Errorf("get: %s err:%s can not get user from db", creators, err)
		h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}
	respondCreators := make([]models.ResponseCreator, len(creators))
	for i, cr := range creators {
		respondCreators[i] = models.ToResponseCreator(cr)
	}

	h.log.Debugf("get creators %s", respondCreators)
	h.Respond(w, r, http.StatusOK, respondCreators)
}
