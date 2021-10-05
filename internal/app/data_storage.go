package app

import (
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/repository"
	"patreon/internal/app/sessions/sessions_manager"
	"patreon/internal/app/store"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type DataStorage struct {
	Store          store.Store
	SessionManager sessions.SessionsManager
}

func (dst *DataStorage) SetStore(store store.Store) {
	dst.Store = store
}
func (dst *DataStorage) SetSessionManager(manager sessions.SessionsManager) {
	dst.SessionManager = manager
}

func NewDataStorage(config *Config, store store.Store) *DataStorage {
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

	sessionManager := sessions_manager.NewSessionManager(repository.NewRedisRepository(redisConn, sessionLog))
	return &DataStorage{
		Store:          store,
		SessionManager: sessionManager,
	}
}
