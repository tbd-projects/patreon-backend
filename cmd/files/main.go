package main

import (
	"flag"
	"os"
	"patreon/internal/app"
	server2 "patreon/internal/microservices/files/delivery/grpc/server"
	repository_os "patreon/internal/microservices/files/files/repository/files/os"
	prometheus_monitoring "patreon/pkg/monitoring/prometheus-monitoring"
	"patreon/pkg/utils"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	grpc2 "google.golang.org/grpc"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}

func main() {
	config := app.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logrus.Fatal(err)
	}
	logger, CloseLogger := utils.NewLogger(config, true, "files_microservice")
	defer CloseLogger()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		os.Exit(1)
	}
	logger.SetLevel(level)

	metrics := prometheus_monitoring.NewPrometheusMetrics("files")
	if err = metrics.SetupMonitoring(); err != nil {
		logger.Fatal(err)
	}

	grpcDurationMetrics := utils.AuthInterceptor(metrics)

	grpc := grpc2.NewServer(
		grpc2.UnaryInterceptor(grpcDurationMetrics),
		grpc2.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)

	grpc_prometheus.Register(grpc)
	filesRepository := repository_os.NewFileRepository(config.MediaDir)
	logger.Info("Files-service create repository")

	server := server2.NewFilesGRPCServer(logger, grpc, filesRepository)
	if err = server.StartGRPCServer(config.Microservices.FilesUrl); err != nil {
		logger.Fatalln(err)
	}
	logger.Info("Files-service was stopped")

}
