package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	_ "patreon/docs"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_factory"
	"patreon/internal/app/store/sqlstore"
	"time"

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
func NewDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
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
	currentTime := time.Now()
	formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + ".out"

	f, err := os.OpenFile(formatted, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	logger.SetOutput(f)
	logger.SetLevel(level)
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	CORSConfigure(router)

	db, err := NewDB(config.DataBaseUrl)
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

	st := sqlstore.New(db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dataStorage := app.NewDataStorage(config, st)

	factory := handler_factory.NewFactory(logger, dataStorage)
	hs := factory.GetHandleUrls()

	for url, h := range *hs {
		h.Connect(router.PathPrefix(url))
	}

	log.Info("starting server")

	return http.ListenAndServe(config.BindAddr, router)
}
