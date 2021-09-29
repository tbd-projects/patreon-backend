package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"
	"patreon/internal/models"

	gh "patreon/internal/app/handlers/general_handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	authMiddleware middleware.SessionMiddleware
	Store          store.Store
	SessionManager sessions.SessionsManager
	gh.RespondHandler
	withHideMethod
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		RespondHandler: gh.RespondHandler{},
		withHideMethod: withHideMethod{gh.NewBaseHandler(logrus.New(), "/profile")},
	}
}

func (h *ProfileHandler) SetStore(store store.Store) {
	h.Store = store
}

func (h *ProfileHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.Log())
}

func (h *ProfileHandler) Join(router *mux.Router) {
	h.baseHandler.AddMethod(gh.GET, h.ServeHTTP)
	h.baseHandler.AddMethod(gh.OPTIONAL, h.ServeHTTP)
	h.baseHandler.AddMiddleware(h.authMiddleware.Check, h.authMiddleware.CorsMiddleware)
	h.baseHandler.Join(router)
}

// Profile
// @Summary get information from user for profile
// @Description get nickname and avatar for profile
// @Accept  json
// @Produce json
// @Success 201 {object} models.ProfileResponse "Successfully get profile"
// @Failure 401 "User are not authorized"
// @Failure 503 {object} models.BaseResponse "Not found user"
// @Failure 500 {object} models.BaseResponse "Error context"
// @Router /profile [GET]
func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log().Error(err)
		}
	}(r.Body)

	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log().Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ContextError)
		return
	}

	u, err := h.Store.User().FindByID(userID.(int64))
	if err != nil {
		h.Log().Errorf("get: %s err:%s can not get user from db", u, err)
		h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}

	h.Log().Debugf("get profile %s", u)
	h.Respond(w, r, http.StatusOK, models.Profile{Nickname: u.Nickname, Avatar: u.Avatar})
}
