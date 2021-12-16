package utilits

import (
	"github.com/mailru/easyjson"
	"net/http"
	"patreon/internal/app/delivery/http/models"
)

type Responder struct {
	LogObject
}

func (h *Responder) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(w, r, code, http_models.ErrResponse{Err: err.Error()})
}

func (h *Responder) Respond(w http.ResponseWriter, r *http.Request, code int, data easyjson.Marshaler) {
	w.WriteHeader(code)
	if data != nil {
		_, _, err := easyjson.MarshalToHTTPResponseWriter(data, w)
		if err != nil {
			h.Log(r).Error(err)
		}
	}
	logUser, _ := easyjson.Marshal(data)
	h.Log(r).Info("Respond data: ", string(logUser))
}
