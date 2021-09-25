package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RespondHandler struct {
	log *logrus.Logger
}

func (h *RespondHandler) PrintRequest(r *http.Request) {
	h.log.Infof("Request: %s. From URL: %s", r.Method, r.URL.Host + r.URL.Path)
}

func (h *RespondHandler) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(w, r, code, map[string]string{"error": err.Error()})
}

func (h *RespondHandler) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(code)
	if data != nil {
		err := encoder.Encode(data)
		if err != nil {
			h.log.Error(err)
		}
	}
	logUser, _ := json.Marshal(data)
	h.log.Info("Respond data: ", string(logUser))
}
