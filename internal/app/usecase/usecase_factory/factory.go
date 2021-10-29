package usecase_factory

import (
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/sessions_manager"
	useAwards "patreon/internal/app/usecase/awards"
	useCreator "patreon/internal/app/usecase/creator"
	useLikes "patreon/internal/app/usecase/likes"
	usePosts "patreon/internal/app/usecase/posts"
	useUser "patreon/internal/app/usecase/user"
)

type UsecaseFactory struct {
	repositoryFactory RepositoryFactory
	userUsecase       useUser.Usecase
	creatorUsecase    useCreator.Usecase
	awardsUsecase     useAwards.Usecase
	sessinManager     sessions.SessionsManager
	postsUsecase      usePosts.Usecase
	likesUsecase      useLikes.Usecase
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

func (f *UsecaseFactory) GetSessionManager() sessions.SessionsManager {
	if f.sessinManager == nil {
		f.sessinManager = sessions_manager.NewSessionManager(f.repositoryFactory.GetSessionRepository())
	}
	return f.sessinManager
}

func (f *UsecaseFactory) GetAwardsUsecase() useAwards.Usecase {
	if f.awardsUsecase == nil {
		f.awardsUsecase = useAwards.NewAwardsUsecase(f.repositoryFactory.GetAwardsRepository())
	}
	return f.awardsUsecase
}

func (f *UsecaseFactory) GetPostsUsecase() usePosts.Usecase {
	if f.postsUsecase == nil {
		f.postsUsecase = usePosts.NewPostsUsecase(f.repositoryFactory.GetPostsRepository(),
			f.repositoryFactory.GetPostsDataRepository())
	}
	return f.postsUsecase
}

func (f *UsecaseFactory) GetLikesUsecase() useLikes.Usecase {
	if f.likesUsecase == nil {
		f.likesUsecase = useLikes.NewLikesUsecase(f.repositoryFactory.GetLikesRepository())
	}
	return f.likesUsecase
}
