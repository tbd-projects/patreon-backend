package base_handler

import (
	"net/http"
	"patreon/internal/app/middleware"
	"strings"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

const (
	GET     = http.MethodGet
	POST    = http.MethodPost
	PUT     = http.MethodPut
	DELETE  = http.MethodDelete
	OPTIONS = http.MethodOptions
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type MiddlewareFunc func(handler http.Handler) http.Handler

type BaseHandler struct {
	utilitiesMiddleware middleware.UtilitiesMiddleware
	handlerMethods      map[string]HandlerFunc
	middlewares         []MiddlewareFunc
	RespondHandler
}

func NewBaseHandler(log *logrus.Logger) *BaseHandler {
	h := &BaseHandler{handlerMethods: map[string]HandlerFunc{}, middlewares: []MiddlewareFunc{},
		utilitiesMiddleware: middleware.NewUtilitiesMiddleware(log)}
	h.log = log
	h.AddMiddleware(h.utilitiesMiddleware.UpgradeLogger, h.utilitiesMiddleware.UpgradeLogger)
	return h
}

func (h *BaseHandler) AddMiddleware(middleware ...MiddlewareFunc) {
	h.middlewares = append(h.middlewares, middleware...)
}

func (h *BaseHandler) AddMethod(method string, handlerMethod HandlerFunc) {
	h.handlerMethods[method] = handlerMethod
}

func (h *BaseHandler) applyMiddleware(handler http.Handler) http.Handler {
	resultHandler := handler
	for _, mw := range h.middlewares {
		resultHandler = mw(resultHandler)
	}
	return resultHandler
}

func (h *BaseHandler) getListMethods() []string {
	var useMethods []string
	for key := range h.handlerMethods {
		useMethods = append(useMethods, key)
	}
	return useMethods
}

func (h *BaseHandler) Connect(route *mux.Route) {
	route.Handler(h.applyMiddleware(h)).Methods(h.getListMethods()...)
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.PrintRequest(r)
	ok := true
	var handler HandlerFunc

	handler, ok = h.handlerMethods[r.Method]
	if ok {
		handler(w, r)
	} else {
		h.log.Errorf("Unexpected http method: %s", r.Method)
		w.Header().Set("Allow", strings.Join(h.getListMethods(), ", "))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
