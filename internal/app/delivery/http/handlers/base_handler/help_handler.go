package base_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	usePosts "patreon/internal/app/usecase/posts"
	"patreon/internal/app/utilits"
	repFiles "patreon/internal/microservices/files/files/repository/files"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

const (
	EmptyQuery      = -2

)

type Sanitizable interface {
	Sanitize(sanitizer bluemonday.Policy)
}

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

func (h *HelpHandlers) redirect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("page", "1")
	query.Set("limit", fmt.Sprintf("%d", usePosts.BaseLimit))
	r.URL.RawQuery = query.Encode()
	redirectUrl := r.URL.String()
	h.Log(r).Debugf("redirect to url: %s, with offest 0 and limit %d", redirectUrl, usePosts.BaseLimit)

	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

// GetPaginationFromQuery Expected api param:
// 	Auto redirect if some param not passed
//	Param page query uint64 true "start page number of posts mutually exclusive with offset"
// 	Param offset query uint64 true "start number of posts mutually exclusive with page"
// 	Param limit query uint64 true "posts to return"
//
// Errors:
// 	Status 400 handler_errors.InvalidQueries
func (h *HelpHandlers) GetPaginationFromQuery(w http.ResponseWriter, r *http.Request) (int64, int64, bool) {
	var limit, offset, page int64
	var ok bool
	limit, ok = h.GetInt64FromQueries(w, r, "limit")
	if !ok {
		if limit == EmptyQuery {
			h.redirect(w, r)
		}
		return app.InvalidInt, app.InvalidInt, false
	}

	offset, ok = h.GetInt64FromQueries(w, r, "offset")
	if !ok {
		if offset != EmptyQuery {
			return app.InvalidInt, app.InvalidInt, false
		}
		page, ok = h.GetInt64FromQueries(w, r, "page")
		if !ok {
			if offset == EmptyQuery {
				h.redirect(w, r)
			}
			return app.InvalidInt, app.InvalidInt, false
		}
		if page <= 0 {
			page = 1
		}
		offset = (page - 1) * limit
	}
	return limit, offset, true
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

// GerFilesFromRequest http Errors:
// 		Status 400 handler_errors.FileSizeError
// 		Status 400 handler_errors.InvalidFormFieldName
// 		Status 400 handler_errors.InvalidImageExt
// 		Status 500 handler_errors.InternalError
func (h *HelpHandlers) GerFilesFromRequest(w http.ResponseWriter, r *http.Request, maxSize int64,
	name string, validTypes []string) (io.Reader, repFiles.FileName, int, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	if err := r.ParseMultipartForm(maxSize); err != nil {
		return nil, "", http.StatusBadRequest, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.FileSizeError,
		}
	}

	f, fHeader, err := r.FormFile(name)
	if err != nil {
		return nil, "", http.StatusBadRequest, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InvalidFormFieldName,
		}
	}

	buff := make([]byte, 512)
	if _, err = f.Read(buff); err != nil {
		return nil, "", http.StatusInternalServerError, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InternalError,
		}
	}
	sort.Strings(validTypes)
	fType := http.DetectContentType(buff)
	if pos := sort.SearchStrings(validTypes, fType); pos == len(validTypes) || validTypes[pos] != fType {
		return nil, "", http.StatusBadRequest, fmt.Errorf("%s, %s",
			handler_errors.InvalidExt, strings.Join(validTypes, " ,"))
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return nil, "", http.StatusInternalServerError, app.GeneralError{
			ExternalErr: err,
			Err:         handler_errors.InternalError,
		}
	}

	return f, repFiles.FileName(fHeader.Filename), 0, nil
}

func (h *HelpHandlers) GetRequestBody(w http.ResponseWriter, r *http.Request,
	reqStruct Sanitizable, sanitizer bluemonday.Policy) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(reqStruct); err != nil {
		return err
	}
	reqStruct.Sanitize(sanitizer)

	return nil
}
