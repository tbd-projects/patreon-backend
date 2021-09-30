package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/handlers/urls"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"
	"patreon/internal/models"

	gh "patreon/internal/app/handlers/general_handlers"
)

type CreatorHandler struct {
	authMiddleware middleware.SessionMiddleware
	Store          store.Store
	SessionManager sessions.SessionsManager
	gh.RespondHandler
	withHideMethod
}

func NewCreatorHandler() *CreatorHandler {
	return &CreatorHandler{
		RespondHandler: gh.RespondHandler{},
		withHideMethod: withHideMethod{gh.NewBaseHandler(logrus.New(), urls.Creator)},
	}
}

func (h *CreatorHandler) SetStore(store store.Store) {
	h.Store = store
}

func (h *CreatorHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.Log())
}

func (h *CreatorHandler) Join(router *mux.Router) {
	h.baseHandler.AddMethod(gh.GET, h.ServeHTTP)
	h.baseHandler.AddMethod(gh.OPTIONAL, h.ServeHTTP)
	h.baseHandler.AddMiddleware(h.authMiddleware.Check)
	h.baseHandler.Join(router)
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
			h.Log().Error(err)
		}
	}(r.Body)

	creators, err := h.Store.Creator().GetCreators()
	if err != nil {
		h.Log().Errorf("get: %s err:%s can not get user from db", creators, err)
		h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}

	respondCreators := make([]models.ResponseCreator, len(creators))
	for i, cr := range creators {
		respondCreators[i] = models.ToResponseCreator(cr)
	}

	h.Log().Debugf("get creators %s", respondCreators)
	h.Respond(w, r, http.StatusOK, respondCreators)
}
