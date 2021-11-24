package server

import (
	"context"
	"fmt"
	"net/http"
	prometheus_monitoring "patreon/pkg/monitoring/prometheus-monitoring"

	"google.golang.org/grpc/connectivity"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	//_ "net/http/pprof"
	"net/url"
	"patreon/internal/app/delivery/http/handler_factory"
	"patreon/internal/app/middleware"
	"patreon/internal/app/repository/repository_factory"
	"patreon/internal/app/usecase/usecase_factory"

	"golang.org/x/crypto/acme/autocert"

	_ "patreon/docs"
	"patreon/internal/app"

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

func (s *Server) checkConnection() error {
	if err := s.connections.SqlConnection.Ping(); err != nil {
		return fmt.Errorf("Can't check connection to sql with error %v ", err)
	}

	s.logger.Info("Success check connection to sql db")

	connAccess, err := s.connections.AccessRedisPool.Dial()
	if err != nil {
		return fmt.Errorf("Can't check connection to redis with error: %v ", err)
	}

	s.logger.Info("Success check connection to redis")

	err = connAccess.Close()
	if err != nil {
		return fmt.Errorf("Can't close connection to redis with error: %v ", err)
	}

	state := s.connections.SessionGrpcConnection.GetState()
	if state != connectivity.Ready {
		return fmt.Errorf("Session connection not ready, status is: %s ", state)
	}

	return nil
}

//return http[0] and https[1] servers
func makingHTTPSServerWithRedirect(config *app.Config, router http.Handler) (*http.Server, *http.Server) {
	serverHTTP := &http.Server{
		Addr: config.BindHttpAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			targetUrl := url.URL{Scheme: "https", Host: r.Host, Path: r.URL.Path, RawQuery: r.URL.RawQuery}
			log.Infof("Redirect from %s, to %s", r.RequestURI, targetUrl.String())
			http.Redirect(w, r, targetUrl.String(), http.StatusPermanentRedirect)
		}),
	}

	hostPolicy := func(ctx context.Context, host string) error {
		allowedHost := config.Domen
		if host == allowedHost {
			return nil
		}
		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	dataDir := "./patreon-secrt"
	m := &autocert.Manager{
		Cache:      autocert.DirCache(dataDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
	}

	serverHTTPS := &http.Server{
		Addr:      config.BindHttpsAddr,
		TLSConfig: m.TLSConfig(),
		Handler:   router,
	}
	return serverHTTP, serverHTTPS
}

// @title Patreon
// @version 1.0
// @description Server for Patreon application.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /api/v1/

// @x-extension-openapi {"example": "value on a json format"}

func (s *Server) Start(config *app.Config) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	monitoringHandler := prometheus_monitoring.NewPrometheusMetrics("main")
	err := monitoringHandler.SetupMonitoring()
	if err != nil {
		return err
	}

	routerApi := router.PathPrefix("/api/v1/").Subrouter()
	routerApi.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//routerApi.HandleFunc("/debug/pprof/", pprof.Index)
	//routerApi.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//routerApi.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//routerApi.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//routerApi.HandleFunc("/debug/pprof/trace", pprof.Trace)

	fileServer := http.FileServer(http.Dir(config.MediaDir + "/"))
	routerApi.PathPrefix("/" + app.LoadFileUrl).Handler(http.StripPrefix("/api/v1/"+app.LoadFileUrl, fileServer))

	repositoryFactory := repository_factory.NewRepositoryFactory(s.logger, s.connections)

	usecaseFactory := usecase_factory.NewUsecaseFactory(repositoryFactory, s.connections.FilesGrpcConnection)
	factory := handler_factory.NewFactory(s.logger, usecaseFactory, s.connections.SessionGrpcConnection)
	hs := factory.GetHandleUrls()

	for apiUrl, h := range *hs {
		h.Connect(routerApi.Path(apiUrl))
	}
	utilitsMiddleware := middleware.NewUtilitiesMiddleware(s.logger, monitoringHandler)
	ddosMiddleware := middleware.NewDdosMiddleware(s.logger, usecaseFactory.GetAccessUsecase())
	routerApi.Use(utilitsMiddleware.CheckPanic, utilitsMiddleware.UpgradeLogger, ddosMiddleware.CheckAccess)

	cors := middleware.NewCorsMiddleware(&config.Cors, router)
	routerCors := cors.SetCors(router)

	if config.IsHTTPSServer {
		serverHTTP, serverHTTPS := makingHTTPSServerWithRedirect(config, routerCors)

		go func(logger *log.Logger, server *http.Server) {
			logger.Info("Start http server")
			err := server.ListenAndServe()
			if err != nil {
				logger.Errorf("http server error on listenAndServe %s", err)
			}
		}(s.logger, serverHTTP)

		s.logger.Info("Start https server")
		return serverHTTPS.ListenAndServeTLS("", "")
	} else {
		s.logger.Info("start no production http server")
		return http.ListenAndServe(config.BindHttpAddr, routerCors)
	}
}
