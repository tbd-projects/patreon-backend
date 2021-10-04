package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions/middleware"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type LogoutHandler struct {
	baseHandler    app.HandlerJoiner
	dataStorage    *app.DataStorage
	authMiddleware middleware.SessionMiddleware
	RespondHandler
}

func NewLogoutHandler(storage *app.DataStorage) *LogoutHandler {
	h := &LogoutHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/logout"),
		dataStorage:    storage,
		RespondHandler: RespondHandler{logrus.New()},
	}
	if storage != nil {
		h.authMiddleware = *middleware.NewSessionMiddleware(h.dataStorage.SessionManager, h.log)
	}
	return h
}
func (h *LogoutHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.Check(h)).Methods("GET")
	h.baseHandler.Join(router)
}

// Profile
// @Summary logout user
// @Description logout user
// @Accept  json
// @Produce json
// @Success 201 {object} models.BaseResponse "Successfully logout"
// @Failure 500 {object} models.BaseResponse "Error logout session"
// @Failure 401 "User not are authorized"
// @Router /logout [GET]
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Fatal(err)
		}
	}(r.Body)
	uniqID := r.Context().Value("uniq_id")
	if uniqID == nil {
		h.log.Error("can not get uniq_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ContextError)
		return
	}

	h.log.Debugf("Logout session: %s", uniqID)

	err := h.dataStorage.SessionManager.Delete(uniqID.(string))
	if err != nil {
		h.log.Errorf("can not delete session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.DeleteCookieFail)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   uniqID.(string),
		Expires: time.Now().AddDate(0, 0, -1),
	}
	http.SetCookie(w, cookie)
	h.Respond(w, r, http.StatusOK, "successfully logout")
}
