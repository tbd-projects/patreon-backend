package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	config  *Config
	handler Handler
	logger  *log.Logger
}

func New(config *Config, handler Handler) *Server {
	return &Server{
		config:  config,
		logger:  log.New(),
		handler: handler,
	}
}

// Start server
func (s *Server) Start() error {
	s.logger.Info("starting server")

	return http.ListenAndServe(s.config.BindAddr, s.handler)
}
