package base_handler

import (
	"encoding/json"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/models"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

type RespondError struct {
	Code  int
	Error error
	Level logrus.Level
}

type CodeMap map[error]RespondError

type RespondHandler struct {
	log   *logrus.Logger
	entry *logrus.Entry
}

func (h *RespondHandler) Log(r *http.Request) *logrus.Entry {
	ctxLogger := r.Context().Value("logger")
	logger := h.log.WithField("urls", r.URL)
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
	h.Respond(w, r, code, models.ErrResponse{Err: err.Error()})
}

func (h *RespondHandler) UsecaseError(w http.ResponseWriter, r *http.Request, usecaseErr error, codeByErr CodeMap) {
	var generalError *app.GeneralError

	if errors.As(usecaseErr, &generalError) {
		usecaseErr = errors.Cause(usecaseErr).(*app.GeneralError).Err
	}

	respond := RespondError{http.StatusServiceUnavailable,
		errors.New("UnknownError"), logrus.ErrorLevel}
	for err, respondErr := range codeByErr {
		if errors.Is(usecaseErr, err) {
			respond = respondErr
			break
		}
	}

	h.Log(r).Logf(respond.Level, "Gotted error: %v", usecaseErr)
	h.Error(w, r, respond.Code, respond.Error)
}
func (h *RespondHandler) HandlerError(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Log(r).Errorf("Gotted error: %v", err)

	var generalError *app.GeneralError
	if errors.As(err, &generalError) {
		err = errors.Cause(err).(*app.GeneralError).Err
	}
	h.Error(w, r, code, err)
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
