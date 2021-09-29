package general_handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	"strings"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type MiddlewareFunc func(handler http.Handler) http.Handler

type BaseHandler struct {
	log            *logrus.Logger
	joinedHandlers *HandlerJoiner
	useMethods     []string
	handlerMethods map[string]HandlerFunc
	middlewares    []MiddlewareFunc
}

func (h *BaseHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}

func NewBaseHandler(log *logrus.Logger, url string) *BaseHandler {
	return &BaseHandler{log: log, joinedHandlers: &HandlerJoiner{currentUrl: url},
		handlerMethods: map[string]HandlerFunc{}, middlewares: []MiddlewareFunc{}, useMethods: []string{}}
}

func (h *BaseHandler) Log() *logrus.Logger {
	return h.log
}

func (h *BaseHandler) AddMethod(method string, handlerMethod HandlerFunc) {
	h.useMethods = append(h.useMethods, method)
	h.handlerMethods[method] = handlerMethod
}

func (h *BaseHandler) AddMiddleware(middleware ...MiddlewareFunc) {
	h.middlewares = middleware
}

func (h *BaseHandler) applyMiddleware(handler http.Handler) http.Handler {
	resultHandler := handler
	for _, mw := range h.middlewares {
		resultHandler = mw(resultHandler)
	}
	return resultHandler
}

func (h *BaseHandler) Join(router *mux.Router) {
	router.Handle(h.joinedHandlers.GetUrl(), h.applyMiddleware(h)).Methods(h.useMethods...)
	h.joinedHandlers.Join(router)
}

func (h *BaseHandler) JoinHandlers(joinedHandlers []app.Joinable) {
	h.joinedHandlers.AddHandlers(joinedHandlers)
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.PrintRequest(r)
	ok := true
	var handler HandlerFunc

	switch r.Method {
	case GET:
		handler, ok = h.handlerMethods[GET]
		handler(w, r)
	case POST:
		handler, ok = h.handlerMethods[POST]
		handler(w, r)
	case OPTIONAL:
		handler, ok = h.handlerMethods[OPTIONAL]
		handler(w, r)
	case PUT:
		handler, ok = h.handlerMethods[PUT]
		handler(w, r)
	default:
		ok = false
	}

	if !ok {
		h.log.Errorf("Unexpected http method: %s", r.Method)
		w.Header().Set("Allow", strings.Join(h.useMethods, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *BaseHandler) PrintRequest(r *http.Request) {
	h.log.Infof("Request: %s. From URL: %s", r.Method, r.URL.Host+r.URL.Path)
}
