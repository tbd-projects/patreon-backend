package usecase_factory

import (
	repCreator "patreon/internal/app/repository/creator"
	repUser "patreon/internal/app/repository/user"
	"patreon/internal/app/sessions"
)

type RepositoryFactory interface {
	GetUserRepository() repUser.Repository
	GetCreatorRepository() repCreator.Repository
	GetSessionRepository() sessions.SessionRepository
}
