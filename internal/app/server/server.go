package server

import (
	"net/http"
	"os"
	_ "patreon/docs"
	"patreon/internal/app"
	"patreon/internal/app/handlers"

	gorilla_handlers "github.com/gorilla/handlers"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config  *app.Config
	handler app.Handler
	logger  *log.Logger
}

func New(config *app.Config, handler app.Handler) *Server {
	return &Server{
		config:  config,
		logger:  log.New(),
		handler: handler,
	}
}
func CORSConfigure(router *mux.Router) {
	if router != nil {
		router.Use(gorilla_handlers.CORS(
			gorilla_handlers.AllowedOrigins([]string{"http://localhost:3001", "https://patreon-dev.herokuapp.com",
				"https://dev-volodya-patreon.netlify.app", "https://patreon.netlify.app",
				"http://patreon-dev.herokuapp.com", "http://front2.tp.volodyalarin.site", "http://pyaterochka-team.site"}),
			gorilla_handlers.AllowedHeaders([]string{
				"Accept", "Content-Type", "Content-Length",
				"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
			gorilla_handlers.AllowCredentials(),
			gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		))
	}
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
func Start(config *app.Config) error {
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

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	CORSConfigure(router)

	handler.SetRouter(router)

	dataStorage := app.NewDataStorage(config)

	registerHandler := handlers.NewRegisterHandler(dataStorage)
	loginHandler := handlers.NewLoginHandler(dataStorage)
	logoutHandler := handlers.NewLogoutHandler(dataStorage)
	profileHandler := handlers.NewProfileHandler(dataStorage)
	creatorHandler := handlers.NewCreatorHandler(dataStorage)
	creatorCreateHandler := handlers.NewCreatorCreateHandler(dataStorage)

	creatorHandler.JoinHandlers([]app.Joinable{
		creatorCreateHandler,
	})

	handler.JoinHandlers([]app.Joinable{
		registerHandler,
		loginHandler,
		profileHandler,
		logoutHandler,
		creatorHandler,
	})

	s := New(config, handler)
	s.logger.Info("starting server")

	return http.ListenAndServe(config.BindAddr, s.handler)
}
