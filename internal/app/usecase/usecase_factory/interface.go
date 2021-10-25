package usecase_factory

import (
	repoAwrds "patreon/internal/app/repository/awards"
	repCreator "patreon/internal/app/repository/creator"
	repUser "patreon/internal/app/repository/user"
	"patreon/internal/app/sessions"
)

type RepositoryFactory interface {
	GetUserRepository() repUser.Repository
	GetCreatorRepository() repCreator.Repository
	GetAwardsRepository() repoAwrds.Repository
	GetSessionRepository() sessions.SessionRepository
}
