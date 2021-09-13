package server

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Server struct {
	config *Config
	logger *log.Logger
	router *mux.Router
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: log.New(),
		router: mux.NewRouter(),
	}
}

// Start server
func (s *Server) Start() error {
	if err := s.configureServer(); err != nil {
		return err
	}
	s.configfureRouter()

	s.logger.Info("starting server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureServer() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *Server) configfureRouter() {
	s.router.HandleFunc("/hello", s.HandleRoot())
}

func (s *Server) HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello patron!")
		if err != nil {
			log.Fatal(err)
		}
	}
}
