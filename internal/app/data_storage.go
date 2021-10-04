package app

import (
	"database/sql"
	"os"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/repository"
	"patreon/internal/app/sessions/sessions_manager"
	"patreon/internal/app/store"
	"patreon/internal/app/store/sqlstore"

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

func NewDataStorage(config *Config) *DataStorage {
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
		Store:          st,
		SessionManager: sessionManager,
	}
}
