package server

import (
	"database/sql"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"patreon/internal/app"
	"patreon/internal/app/handlers"
	"patreon/internal/app/sessions/repository"
	"patreon/internal/app/sessions/sessions_manager"
	"patreon/internal/app/store/sqlstore"

	gorilla_handlers "github.com/gorilla/handlers"

	redis "github.com/gomodule/redigo/redis"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	router.Use(gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"http://localhost:3001", "https://patreon-dev.herokuapp.com",
			"https://dev-volodya-patreon.netlify.app", "https://patreon.netlify.app",
			"http://patreon-dev.herokuapp.com"}),
		gorilla_handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "csrf-token", "Authorization"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	))

	handler.SetRouter(router)

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

	registerHandler := handlers.NewRegisterHandler()
	registerHandler.SetStore(store)

	loginHandler := handlers.NewLoginHandler()
	loginHandler.SetStore(store)

	profileHandler := handlers.NewProfileHandler()
	profileHandler.SetStore(store)

	logoutHandler := handlers.NewLogoutHandler()

	sessionLog := log.New()
	sessionLog.SetLevel(log.FatalLevel)
	redisConn := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.RedisUrl)
		},
	}

	conn, err := redisConn.Dial()
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	redisRepository := repository.NewRedisRepository(redisConn, sessionLog)
	sessionManager := sessions_manager.NewSessionManager(redisRepository)
	loginHandler.SetSessionManager(sessionManager)
	registerHandler.SetSessionManager(sessionManager)
	profileHandler.SetSessionManager(sessionManager)
	logoutHandler.SetSessionManager(sessionManager)

	handler.JoinHandlers([]app.Joinable{
		registerHandler,
		loginHandler,
		profileHandler,
		logoutHandler,
	})

	s := New(config, handler)

	m := &autocert.Manager{
		Cache:      autocert.DirCache("golang-autocert"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Domen),
	}
	serverHTTPS := &http.Server{
		Addr:      config.BindAddrHTTPS,
		TLSConfig: m.TLSConfig(),
		Handler:   s.handler,
	}

	serverHTTP := &http.Server{
		Addr: config.BindAddrHTTP,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+config.Domen+r.RequestURI, http.StatusMovedPermanently)
		}),
	}

	s.logger.Info("starting server")

	go func(log *log.Logger) {
		err := serverHTTP.ListenAndServe()
		if err != nil {
			log.Fatalf("Drop http server with error: %s", err)
		}
		log.Info("starting http server")
	}(s.logger)

	return serverHTTPS.ListenAndServeTLS("", "")
}
