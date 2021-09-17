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
	if err := s.configureServer(); err != nil {
		return err
	}

	s.logger.Info("starting server")

	return http.ListenAndServe(s.config.BindAddr, s.handler)
}
func (s *Server) configureServer() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}
