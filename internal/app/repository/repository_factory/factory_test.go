package repository_factory

import (
	"patreon/internal/app"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	t.Helper()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	redisServer, err := miniredis.Run()
	require.NoError(t, err)

	addr := redisServer.Addr()
	redisConn := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}

	factory := NewRepositoryFactory(log, app.ExpectedConnections{SqlConnection: db, AccessRedisPool: redisConn, PathFiles: "don/"})
	factory.GetFilesRepository()
	factory.GetPostsDataRepository()
	factory.GetLikesRepository()
	factory.GetAwardsRepository()
	factory.GetCreatorRepository()
	factory.GetUserRepository()
	factory.GetAccessRepository()
	factory.GetCsrfRepository()
	factory.GetPostsRepository()
	factory.GetSubscribersRepository()
	factory.GetInfoRepository()
	factory.GetPaymentsRepository()

	redisServer.Close()
}
