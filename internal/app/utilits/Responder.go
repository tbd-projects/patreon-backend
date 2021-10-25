package utilits

import (
	"encoding/json"
	"net/http"
	"patreon/internal/app/delivery/http/models"
)

type Responder struct {
	LogObject
}

func (h *Responder) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(w, r, code, models.ErrResponse{Err: err.Error()})
}

func (h *Responder) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(code)
	if data != nil {
		err := encoder.Encode(data)
		if err != nil {
			h.Log(r).Error(err)
		}
	}
	logUser, _ := json.Marshal(data)
	h.Log(r).Info("Respond data: ", string(logUser))
}
