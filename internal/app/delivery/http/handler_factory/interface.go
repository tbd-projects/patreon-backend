package handler_factory

import (
	useCsrf "patreon/internal/app/csrf/usecase"
	"patreon/internal/app/sessions"
	useAwards "patreon/internal/app/usecase/awards"
	useCreator "patreon/internal/app/usecase/creator"
	useSubscr "patreon/internal/app/usecase/subscribers"
	useUser "patreon/internal/app/usecase/user"
)

type UsecaseFactory interface {
	GetUserUsecase() useUser.Usecase
	GetCreatorUsecase() useCreator.Usecase
	GetCsrfUsecase() useCsrf.Usecase
	GetAwardsUsecase() useAwards.Usecase
	GetSessionManager() sessions.SessionsManager
	GetSubscribersUsecase() useSubscr.Usecase
}
