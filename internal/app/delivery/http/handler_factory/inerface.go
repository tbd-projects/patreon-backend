package handler_factory

import (
	"patreon/internal/app/sessions"
	useAwards "patreon/internal/app/usecase/awards"
	useCreator "patreon/internal/app/usecase/creator"
	useUser "patreon/internal/app/usecase/user"
)

type UsecaseFactory interface {
	GetUserUsecase() useUser.Usecase
	GetCreatorUsecase() useCreator.Usecase
	GetAwardsUsecase() useAwards.Usecase
	GetSessionManager() sessions.SessionsManager
}
