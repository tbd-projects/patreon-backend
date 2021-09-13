package server

import (
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config *Config
	logger *log.Logger
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: log.New(),
	}
}

// Start server
func (s *Server) Start() error {
	if err := s.configureServer(); err != nil {
		return err
	}

	s.logger.Info("starting server")

	return nil
}

func (s *Server) configureServer() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}
