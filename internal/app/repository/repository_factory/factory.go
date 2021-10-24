package repository_factory

import (
	"github.com/sirupsen/logrus"
	"patreon/internal/app"
	repCreator "patreon/internal/app/repository/creator"
	repCreatorPsql "patreon/internal/app/repository/creator/postgresql"
	repUser "patreon/internal/app/repository/user"
	repUserPsql "patreon/internal/app/repository/user/postgresql"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/repository"
)

type RepositoryFactory struct {
	expectedConnections app.ExpectedConnections
	logger              *logrus.Logger
	userRepository      repUser.Repository
	creatorRepository   repCreator.Repository
	sessinRepository    sessions.SessionRepository
}

func NewRepositoryFactory(logger *logrus.Logger, expectedConnections app.ExpectedConnections) *RepositoryFactory {
	return &RepositoryFactory{
		expectedConnections: expectedConnections,
		logger:              logger,
	}
}

func (f *RepositoryFactory) GetUserRepository() repUser.Repository {
	if f.userRepository == nil {
		f.userRepository = repUserPsql.NewUserRepository(f.expectedConnections.SqlConnection)
	}
	return f.userRepository
}

func (f *RepositoryFactory) GetCreatorRepository() repCreator.Repository {
	if f.creatorRepository == nil {
		f.creatorRepository = repCreatorPsql.NewCreatorRepository(f.expectedConnections.SqlConnection)
	}
	return f.creatorRepository
}

func (f *RepositoryFactory) GetSessionRepository() sessions.SessionRepository {
	if f.sessinRepository == nil {
		f.sessinRepository = repository.NewRedisRepository(f.expectedConnections.RedisPool, f.logger)
	}
	return f.sessinRepository
}
