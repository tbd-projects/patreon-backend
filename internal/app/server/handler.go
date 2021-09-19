package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"patreon/internal/app/store"
	"patreon/internal/models"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type MainHandler struct {
	router *mux.Router
	Store  store.Store
	log    *logrus.Logger
}

func NewMainHandler() *MainHandler {
	return &MainHandler{
		log: logrus.New(),
	}
}

func (h *MainHandler) SetRouter(router *mux.Router) {
	h.router = router
}
func (h *MainHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *MainHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *MainHandler) HandleRegistration() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.log.Error(err)
			}
		}(r.Body)
		req := &request{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(req); err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u := &models.User{
			Login:    req.Login,
			Password: req.Password,
		}
		checkUser, _ := h.Store.User().FindByLogin(u.Login)
		if checkUser != nil {
			h.error(w, r, http.StatusConflict, store.UserAlreadyExist)
			return
		}
		if err := h.Store.User().Create(u); err != nil {
			h.error(w, r, http.StatusCreated, err)
			return
		}
		u.MakePrivateDate()
		h.respond(w, r, http.StatusOK, u)

	}
}

func (h *MainHandler) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.respond(w, r, code, map[string]string{"error": err.Error()})
}
func (h *MainHandler) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(code)
	if data != nil {
		err := encoder.Encode(data)
		if err != nil {
			h.log.Error(err)
		}
	}
	logrus.Info("respond data: ", data)
}

func (h *MainHandler) RegisterHandlers() {
	h.router.HandleFunc("/register", h.HandleRegistration()).Methods("POST")

}
