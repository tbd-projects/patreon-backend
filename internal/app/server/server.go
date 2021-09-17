package server

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"patreon/internal/app/store"
)

type Server struct {
	config *Config
	logger *log.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config, router *mux.Router) *Server {
	return &Server{
		config: config,
		logger: log.New(),
		router: router,
	}
}

// Start server
func (s *Server) Start() error {
	if err := s.configureServer(); err != nil {
		return err
	}
	s.configfureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}
func (s *Server) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

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

func (s *Server) configfureRouter() {
	s.router.HandleFunc("/hello", s.HandleRoot()).Methods("GET", "POST")
}

func (s *Server) HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello patron!")
		if err != nil {
			log.Fatal(err)
		}
	}
}
