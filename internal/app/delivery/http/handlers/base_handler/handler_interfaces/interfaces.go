package handler_interfaces

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request)
type HMiddlewareFunc func(http.Handler) http.Handler
type HFMiddlewareFunc func(HandlerFunc) HandlerFunc
