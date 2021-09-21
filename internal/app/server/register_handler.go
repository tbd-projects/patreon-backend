package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"patreon/internal/app/server/joinable_handler"
	"patreon/internal/app/store"
)

type RegisterHandler struct {
	baseHandler joinable_handler.HandlerJoiner
	Store       store.Store
	log         *logrus.Logger
}

func NewRegisterHandler(store store.Store, joinedHandlers []joinable_handler.IJoinable) *RegisterHandler {
	return &RegisterHandler{
		baseHandler: joinable_handler.CreateHandlerJoiner(joinedHandlers, "/register"),
		Store:       store,
		log:         logrus.New(),
	}
}

func (h *RegisterHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *RegisterHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Error(err)
		}
	}(r.Body)
	//req := &request{}
	/*
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
		h.respond(w, r, http.StatusOK, u)*/
}

func (h *RegisterHandler) Join(router *mux.Router) {
	router.HandleFunc(h.baseHandler.GetUrl(), h.ServeHTTP).Methods("POST")
	h.baseHandler.Join(router)
}

func (h *RegisterHandler) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (h *RegisterHandler) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
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
