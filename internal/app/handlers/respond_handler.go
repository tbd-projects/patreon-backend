package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RespondHandler struct{}

func (h *RespondHandler) Error(log *logrus.Logger, w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(log, w, r, code, map[string]string{"error": err.Error()})
}
func (h *RespondHandler) Respond(log *logrus.Logger, w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(code)
	if data != nil {
		err := encoder.Encode(data)
		if err != nil {
			log.Error(err)
		}
	}
	logUser, _ := json.Marshal(data)
	log.Info("Respond data: ", string(logUser))
}
