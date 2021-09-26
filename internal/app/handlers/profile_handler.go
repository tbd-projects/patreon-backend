package handlers

import (
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	baseHandler    app.HandlerJoiner
	authMiddleware middleware.SessionMiddleware
	router         *mux.Router
	Store          store.Store
	SessionManager sessions.SessionsManager
	RespondHandler
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/profile"),
		RespondHandler: RespondHandler{logrus.New()},
	}
}

func (h *ProfileHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *ProfileHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h *ProfileHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.log)
}
func (h *ProfileHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.CorsMiddleware(h.authMiddleware.Check(h))).Methods("GET", "OPTIONS")
	//router.Use(gorilla_handlers.CORS(
	//	gorilla_handlers.AllowedOrigins([]string{"http://localhost:3001"}),
	//	gorilla_handlers.AllowedHeaders([]string{"*"}),
	//	gorilla_handlers.AllowCredentials(),
	//	gorilla_handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
	//))
	//router.HandleFunc(h.baseHandler.GetUrl(), h.ServeHTTP).Methods("POST", "GET")
	//router.Use(h.authMiddleware.Check)
	h.baseHandler.Join(router)
}
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

	u, err := h.Store.User().FindByID(userID.(int64))
	if err != nil {
		h.log.Errorf("get: %s err:%s can not get user from db", u, err)
		h.Error(w, r, http.StatusServiceUnavailable, handler_errors.GetProfileFail)
		return
	}

	h.log.Debugf("get profile %s", u)
	h.Respond(w, r, http.StatusOK, u)
}
