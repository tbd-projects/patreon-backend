package repository_factory

import (
	"github.com/sirupsen/logrus"
	"patreon/internal/app"
	repoAwrds "patreon/internal/app/repository/awards"
	repAwardsPsql "patreon/internal/app/repository/awards/postgresql"
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
	sessionRepository   sessions.SessionRepository
	awardsRepository    repoAwrds.Repository
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
	if f.sessionRepository == nil {
		f.sessionRepository = repository.NewRedisRepository(f.expectedConnections.RedisPool, f.logger)
	}
	return f.sessionRepository
}

func (f *RepositoryFactory) GetAwardsRepository() repoAwrds.Repository {
	if f.awardsRepository == nil {
		f.awardsRepository = repAwardsPsql.NewAwardsRepository(f.expectedConnections.SqlConnection)
	}
	return f.awardsRepository
}
