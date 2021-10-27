package repository_factory

import (
	"patreon/internal/app"
	repCsrf "patreon/internal/app/csrf/repository/jwt"
	repository_access "patreon/internal/app/repository/access"
	repCreator "patreon/internal/app/repository/creator"
	repUser "patreon/internal/app/repository/user"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/repository"

	"github.com/sirupsen/logrus"
)

type RepositoryFactory struct {
	expectedConnections app.ExpectedConnections
	logger              *logrus.Logger
	userRepository      repUser.Repository
	creatorRepository   repCreator.Repository
	csrfRepository      repCsrf.Repository
	sessionRepository   sessions.SessionRepository
	accessRepository    repository_access.Repository
}

func NewRepositoryFactory(logger *logrus.Logger, expectedConnections app.ExpectedConnections) *RepositoryFactory {
	return &RepositoryFactory{
		expectedConnections: expectedConnections,
		logger:              logger,
	}
}

func (f *RepositoryFactory) GetUserRepository() repUser.Repository {
	if f.userRepository == nil {
		f.userRepository = repUser.NewUserRepository(f.expectedConnections.SqlConnection)
	}
	return f.userRepository
}

func (f *RepositoryFactory) GetCreatorRepository() repCreator.Repository {
	if f.creatorRepository == nil {
		f.creatorRepository = repCreator.NewCreatorRepository(f.expectedConnections.SqlConnection)
	}
	return f.creatorRepository
}

func (f *RepositoryFactory) GetCsrfRepository() repCsrf.Repository {
	if f.csrfRepository == nil {
		f.csrfRepository = repCsrf.NewJwtRepository()
	}
	return f.csrfRepository
}

func (f *RepositoryFactory) GetSessionRepository() sessions.SessionRepository {
	if f.sessionRepository == nil {
		f.sessionRepository = repository.NewRedisRepository(f.expectedConnections.RedisPool, f.logger)
	}
	return f.sessionRepository
}
func (f *RepositoryFactory) GetAccessRepository() repository_access.Repository {
	if f.accessRepository == nil {
		f.accessRepository = repository_access.NewRedisRepository(f.expectedConnections.RedisPool, f.logger)
	}
	return f.accessRepository
}
