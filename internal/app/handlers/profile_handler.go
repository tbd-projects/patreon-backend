package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	baseHandler    app.HandlerJoiner
	dataStorage    *app.DataStorage
	authMiddleware middleware.SessionMiddleware
	RespondHandler
}

func NewProfileHandler(storage *app.DataStorage) *ProfileHandler {
	h := &ProfileHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/profile"),
		RespondHandler: RespondHandler{logrus.New()},
		dataStorage:    storage,
	}
	if storage != nil {
		h.authMiddleware = *middleware.NewSessionMiddleware(h.dataStorage.SessionManager, h.log)
	}
	return h
}

func (h *ProfileHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.Check(h)).Methods("GET", "OPTIONS")
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
			h.log.Error(err)
		}
	}(r.Body)

	userID := r.Context().Value("user_id")
	if userID == nil {
		h.log.Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ContextError)
		return
	}

	u, err := h.dataStorage.Store.User().FindByID(userID.(int64))
	if err != nil {
		h.log.Errorf("get: %s err:%s can not get user from db", u, err)
		h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}

	h.log.Debugf("get profile %s", u)
	h.Respond(w, r, http.StatusOK, models.Profile{Nickname: u.Nickname, Avatar: u.Avatar})
}
