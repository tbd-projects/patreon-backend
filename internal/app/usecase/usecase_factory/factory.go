package usecase_factory

import (
	usecase_csrf "patreon/internal/app/csrf/usecase"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/sessions_manager"
	useAccess "patreon/internal/app/usecase/access"
	useCreator "patreon/internal/app/usecase/creator"
	useUser "patreon/internal/app/usecase/user"
)

type UsecaseFactory struct {
	repositoryFactory RepositoryFactory
	userUsecase       useUser.Usecase
	creatorUsecase    useCreator.Usecase
	csrfUsecase       usecase_csrf.Usecase
	sessionsManager   sessions.SessionsManager
	accessUsecase     useAccess.Usecase
}

func NewUsecaseFactory(repositoryFactory RepositoryFactory) *UsecaseFactory {
	return &UsecaseFactory{
		repositoryFactory: repositoryFactory,
	}
}

func (f *UsecaseFactory) GetUserUsecase() useUser.Usecase {
	if f.userUsecase == nil {
		f.userUsecase = useUser.NewUserUsecase(f.repositoryFactory.GetUserRepository())
	}
	return f.userUsecase
}

func (f *UsecaseFactory) GetCreatorUsecase() useCreator.Usecase {
	if f.creatorUsecase == nil {
		f.creatorUsecase = useCreator.NewCreatorUsecase(f.repositoryFactory.GetCreatorRepository())
	}
	return f.creatorUsecase
}
func (f *UsecaseFactory) GetCsrfUsecase() usecase_csrf.Usecase {
	if f.csrfUsecase == nil {
		f.csrfUsecase = usecase_csrf.NewCsrfUsecase(f.repositoryFactory.GetCsrfRepository())
	}
	return f.csrfUsecase
}

func (f *UsecaseFactory) GetSessionManager() sessions.SessionsManager {
	if f.sessionsManager == nil {
		f.sessionsManager = sessions_manager.NewSessionManager(f.repositoryFactory.GetSessionRepository())
	}
	return f.sessionsManager
}
func (f *UsecaseFactory) GetAccessUsecase() useAccess.Usecase {
	if f.accessUsecase == nil {
		f.accessUsecase = useAccess.NewAccessUsecase(f.repositoryFactory.GetAccessRepository())
	}
	return f.accessUsecase
}
