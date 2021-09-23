package server

import (
	"database/sql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"patreon/internal/app"
	"patreon/internal/app/handlers"
	"patreon/internal/app/store/sqlstore"
)

type Server struct {
	config  *Config
	handler app.Handler
	logger  *log.Logger
}

func New(config *Config, handler app.Handler) *Server {
	return &Server{
		config:  config,
		logger:  log.New(),
		handler: handler,
	}
}

func newDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Start(config *Config) error {
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	logger := log.New()
	logger.SetLevel(level)

	handler := handlers.NewMainHandler()
	handler.SetLogger(logger)

	router := mux.NewRouter()
	handler.SetRouter(router)

	//handler.RegisterHandlers()

	db, err := newDB(config.DataBaseUrl)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	store := sqlstore.New(db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	handler.SetStore(store)

	registerHandler := handlers.NewRegisterHandler()
	registerHandler.SetStore(store)

	loginHandler := handlers.NewLoginHandler()
	loginHandler.SetStore(store)
	//joinedHandlers := []app.Joinable{
	//	handlers.NewRegisterHandler(),
	//}

	handler.JoinHandlers([]app.Joinable{registerHandler, loginHandler})

	s := New(config, handler)
	s.logger.Info("starting server")

	return http.ListenAndServe(config.BindAddr, s.handler)
}
