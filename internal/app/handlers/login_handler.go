package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/store"
)

type LoginHandler struct {
	baseHandler app.HandlerJoiner
	router      *mux.Router
	Store       store.Store
	log         *logrus.Logger
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		baseHandler: *app.NewHandlerJoiner([]app.Joinable{}, "/login"),
		log:         logrus.New(),
	}
}

func (h *LoginHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *LoginHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h *LoginHandler) Join(router *mux.Router) {
	router.HandleFunc(h.baseHandler.GetUrl(), h.ServeHTTP).Methods("POST", "GET")
	h.baseHandler.Join(router)
}
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r.Body)

	req := &request{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(req); err != nil {
		h.Error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	u, err := h.Store.User().FindByLogin(req.Login)

	if err != nil || !u.ComparePassword(req.Password) {
		h.Error(w, r, http.StatusUnauthorized, store.IncorrectEmailOrPassword)
		return
	}

	h.Respond(w, r, http.StatusOK, "successfully login")
}
func (h *LoginHandler) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(w, r, code, map[string]string{"error": err.Error()})
}
func (h *LoginHandler) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(code)
	if data != nil {
		err := encoder.Encode(data)
		if err != nil {
			h.log.Error(err)
		}
	}
	logUser, _ := json.Marshal(data)
	logrus.Info("Respond data: ", string(logUser))
}
