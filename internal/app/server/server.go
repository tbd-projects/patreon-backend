package server

import (
	"fmt"
	gh "github.com/gorilla/handlers"
	"net/http"
	ds "patreon/internal/app/server/data_storage"

	_ "patreon/docs"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_factory"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config      *app.Config
	logger      *log.Logger
	connections app.ExpectedConnections
}

func New(config *app.Config, connections app.ExpectedConnections, logger *log.Logger) *Server {
	return &Server{
		config:      config,
		logger:      logger,
		connections: connections,
	}
}

func CORSConfigure(router *mux.Router) {
	if router != nil {
		router.Use(gh.CORS(
			gh.AllowedOrigins([]string{"http://localhost:3001", "https://patreon-dev.herokuapp.com",
				"https://dev-volodya-patreon.netlify.app", "https://patreon.netlify.app",
				"http://patreon-dev.herokuapp.com", "http://front2.tp.volodyalarin.site", "http://pyaterochka-team.site"}),
			gh.AllowedHeaders([]string{
				"Accept", "Content-Type", "Content-Length",
				"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
			gh.AllowCredentials(),
			gh.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		))
	}
}

func (s *Server) checkConnection() error {
	if err := s.connections.SqlConnection.Ping(); err != nil {
		return fmt.Errorf("Can't check connection to sql with error %v ", err)
	}

	conn, err := s.connections.RedisPool.Dial()
	if err != nil {
		return fmt.Errorf("Can't check connection to redis with error: %v ", err)
	}

	err = conn.Close()
	if err != nil {
		return fmt.Errorf("Can't close connection to redis with error: %v ", err)
	}

	return nil
}

// @title Patreon
// @version 1.0
// @description Server for Patreon application.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}
func (s *Server) Start(config *app.Config) error {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	CORSConfigure(router)

	dataStorage := ds.NewDataStorage(s.connections, s.logger)

	factory := handler_factory.NewFactory(s.logger, dataStorage)
	hs := factory.GetHandleUrls()

	for url, h := range *hs {
		h.Connect(router.PathPrefix(url))
	}

	s.logger.Info("starting server")

	return http.ListenAndServe(config.BindAddr, router)
}
