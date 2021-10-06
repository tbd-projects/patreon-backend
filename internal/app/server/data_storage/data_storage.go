package data_storage

import (
	"github.com/sirupsen/logrus"
	"patreon/internal/app"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/repository"
	sm "patreon/internal/app/sessions/sessions_manager"
	"patreon/internal/app/store"
	"patreon/internal/app/store/sqlstore"
)

type DataStorage struct {
	store          store.Store
	sessionManager sessions.SessionsManager
}

func (dst *DataStorage) Store() store.Store {
	return dst.store
}

func (dst *DataStorage) SessionManager() sessions.SessionsManager {
	return dst.sessionManager
}

func (dst *DataStorage) SetStore(store store.Store) {
	dst.store = store
}

func (dst *DataStorage) SetSessionManager(manager sessions.SessionsManager) {
	dst.sessionManager = manager
}

func NewDataStorage(connections app.ExpectedConnections, log *logrus.Logger) *DataStorage {
	return &DataStorage{
		store:          sqlstore.New(connections.SqlConnection),
		sessionManager: sm.NewSessionManager(repository.NewRedisRepository(connections.RedisPool, log)),
	}
}
