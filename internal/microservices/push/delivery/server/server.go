package push_server

import (
	"context"
	"fmt"
	"net/http"
	"patreon/internal/microservices/auth/delivery/grpc/client"
	"patreon/internal/microservices/push"
	prometheus_monitoring "patreon/pkg/monitoring/prometheus-monitoring"

	"google.golang.org/grpc/connectivity"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"patreon/internal/app/middleware"

	//_ "patreon/docs"
	"patreon/internal/app"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config      *push.Config
	logger      *log.Logger
	connections app.ExpectedConnections
}

func New(config *push.Config, connections app.ExpectedConnections, logger *log.Logger) *Server {
	return &Server{
		config:      config,
		logger:      logger,
		connections: connections,
	}
}

func (s *Server) checkConnection() error {
	s.connections.SessionGrpcConnection.WaitForStateChange(context.Background(), connectivity.Ready)
	state := s.connections.SessionGrpcConnection.GetState()
	if state != connectivity.Ready {
		return fmt.Errorf("Session connection not ready, status is: %s ", state)
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
// @BasePath /api/v1/

// @x-extension-openapi {"example": "value on a json format"}

func (s *Server) Start() error {
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

	h := NewPushHandler(s.logger, client.NewSessionClient(s.connections.SessionGrpcConnection))
	h.Connect(routerApi.Path("/user/push"))

	utilitsMiddleware := middleware.NewUtilitiesMiddleware(s.logger, monitoringHandler)
	routerApi.Use(utilitsMiddleware.CheckPanic, utilitsMiddleware.UpgradeLogger)

	cors := middleware.NewCorsMiddleware(&s.config.Cors, router)
	routerCors := cors.SetCors(router)

	s.logger.Info("start no production http server")
	return http.ListenAndServe(s.config.BindHttpAddr, routerCors)
}
