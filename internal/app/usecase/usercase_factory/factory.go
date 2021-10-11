package usercase_factory

import (
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/sessions_manager"
	useCreator "patreon/internal/app/usecase/creator"
	useUser "patreon/internal/app/usecase/user"
)

type UsecaseFactory struct {
	repositoryFactory RepositoryFactory
	userUsercase      useUser.Usecase
	creatorUsercase   useCreator.Usecase
	sessinManager     sessions.SessionsManager
}

func NewUsecaseFactory(repositoryFactory RepositoryFactory) *UsecaseFactory {
	return &UsecaseFactory{
		repositoryFactory: repositoryFactory,
	}
}

func (f *UsecaseFactory) GetUserUsecase() useUser.Usecase {
	if f.userUsercase == nil {
		f.userUsercase = useUser.NewUserUsecase(f.repositoryFactory.GetUserRepository())
	}
	return f.userUsercase
}

func (f *UsecaseFactory) GetCreatorUsecase() useCreator.Usecase {
	if f.creatorUsercase == nil {
		f.creatorUsercase = useCreator.NewCreatorUsecase(f.repositoryFactory.GetCreatorRepository())
	}
	return f.creatorUsercase
}

func (f *UsecaseFactory) GetSessionManager() sessions.SessionsManager {
	if f.sessinManager == nil {
		f.sessinManager = sessions_manager.NewSessionManager(f.repositoryFactory.GetSessionRepository())
	}
	return f.sessinManager
}
