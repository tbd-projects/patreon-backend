package utils

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"patreon/internal/app"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewLogger(config *app.Config, isService bool, serviceName string) (log *logrus.Logger, closeResource func() error) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	currentTime := time.Now().In(time.UTC)
	var servicePath string
	if isService {
		servicePath = serviceName
	}
	formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + "__" + servicePath + ".log"

	f, err := os.OpenFile(formatted, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logger.SetOutput(f)
	logger.Writer()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, f.Close
}

func NewPostgresConnection(config *app.RepositoryConnections) (db *sqlx.DB, closeResource func() error) {
	db, err := sqlx.Open("postgres", config.DataBaseUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	return db, db.Close
}

func NewRedisPool(redisUrl string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}
}

func NewGrpcConnection(grpcUrl string) (*grpc.ClientConn, error) {
	return grpc.Dial(grpcUrl, grpc.WithInsecure())
}

func StringsToLowerCase(array[]string) []string {
	res := make([]string, len(array))
	for i, str := range array {
		res[i] = strings.ToLower(str)
	}
	return res
}