package base_handler

import (
	"encoding/json"
	"net/http"
	"patreon/internal/app"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

type RespondError struct {
	Code  int
	Error error
}

type CodeMap map[error]RespondError

type RespondHandler struct {
	log   *logrus.Logger
	entry *logrus.Entry
}

func (h *RespondHandler) Log(r *http.Request) *logrus.Entry {
	if h.entry != nil {
		return h.entry
	}
	ctxLogger := r.Context().Value("logger")
	logger := h.log.WithContext(r.Context())
	if ctxLogger != nil {
		if log, ok := ctxLogger.(*logrus.Entry); ok {
			logger = log
		}
	}
	h.entry = logger
	return h.entry
}

func (h *RespondHandler) PrintRequest(r *http.Request) {
	h.Log(r).Infof("Request: %s. From URL: %s", r.Method, r.URL.Host+r.URL.Path)
}

func (h *RespondHandler) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(w, r, code, map[string]string{"error": err.Error()})
}

func (h *RespondHandler) UsecaseError(w http.ResponseWriter, r *http.Request, usecaseErr error, codeByErr CodeMap) {
	h.Log(r).Errorf("Gotted error: %v", usecaseErr)

	var generalError *app.GeneralError
	if errors.As(usecaseErr, &generalError) {
		usecaseErr = errors.Cause(usecaseErr).(*app.GeneralError).Err
	}

	for err, respondErr := range codeByErr {
		if errors.Is(usecaseErr, err) {
			h.Error(w, r, respondErr.Code, respondErr.Error)
			return
		}
	}
	h.Error(w, r, http.StatusServiceUnavailable, errors.New("UnknownError"))
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
	h.Log(r).Info("Respond data: ", string(logUser))
}
