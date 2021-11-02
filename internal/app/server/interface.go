package server

import (
	"patreon/internal/app"
)

type HandlerFactory interface {
	GetHandleUrls() *map[string]app.Handler
}
