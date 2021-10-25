package handler_factory

import (
	useCsrf "patreon/internal/app/csrf/usecase"
	"patreon/internal/app/sessions"
	useCreator "patreon/internal/app/usecase/creator"
	useUser "patreon/internal/app/usecase/user"
)

type UsecaseFactory interface {
	GetUserUsecase() useUser.Usecase
	GetCreatorUsecase() useCreator.Usecase
	GetCsrfUsecase() useCsrf.Usecase
	GetSessionManager() sessions.SessionsManager
}
