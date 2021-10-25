package usecase_factory

import (
	repCsrf "patreon/internal/app/csrf/repository/jwt"
	repCreator "patreon/internal/app/repository/creator"
	repUser "patreon/internal/app/repository/user"

	"patreon/internal/app/sessions"
)

type RepositoryFactory interface {
	GetUserRepository() repUser.Repository
	GetCreatorRepository() repCreator.Repository
	GetCsrfRepository() repCsrf.Repository
	GetSessionRepository() sessions.SessionRepository
}
