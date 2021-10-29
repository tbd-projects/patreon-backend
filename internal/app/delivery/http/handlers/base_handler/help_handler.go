package base_handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/utilits"
	"strconv"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

const EmptyQuery = -2

type RespondError struct {
	Code  int
	Error error
	Level logrus.Level
}

type CodeMap map[error]RespondError

type HelpHandlers struct {
	utilits.Responder
}

func (h *HelpHandlers) PrintRequest(r *http.Request) {
	h.Log(r).Infof("Request: %s. From URL: %s", r.Method, r.URL.Host+r.URL.Path)
}

// GetInt64FromParam HTTPErrors
//		Status 400 handler_errors.InvalidParameters
func (h *HelpHandlers) GetInt64FromParam(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	vars := mux.Vars(r)
	number, ok := vars[name]
	numberInt, err := strconv.ParseInt(number, 10, 64)
	if !ok || err != nil {
		h.Log(r).Infof("can'not get parametrs %s, was got %v)", name, vars)
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return app.InvalidInt, false
	}
	return numberInt, true
}

// GetInt64FromQueries HTTPErrors
//		Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetInt64FromQueries(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	number := r.URL.Query().Get(name)
	if number == "" {
		return EmptyQuery, false
	}

	numberInt, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		h.Log(r).Infof("can'not get parametrs %s from query url %s)", name, r.URL)
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidQueries)
		return app.InvalidInt, false
	}
	return numberInt, true
}

func (h *HelpHandlers) UsecaseError(w http.ResponseWriter, r *http.Request, usecaseErr error, codeByErr CodeMap) {
	var generalError *app.GeneralError
	orginalError := usecaseErr
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

	h.Log(r).Logf(respond.Level, "Gotted error: %v", orginalError)
	h.Error(w, r, respond.Code, respond.Error)
}

func (h *HelpHandlers) HandlerError(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Log(r).Errorf("Gotted error: %v", err)

	var generalError *app.GeneralError
	if errors.As(err, &generalError) {
		err = errors.Cause(err).(*app.GeneralError).Err
	}
	h.Error(w, r, code, err)
}
