package usecase_factory

import (
	usecase_csrf "patreon/internal/app/csrf/usecase"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/sessions_manager"
	useAccess "patreon/internal/app/usecase/access"
	useAwards "patreon/internal/app/usecase/awards"
	useCreator "patreon/internal/app/usecase/creator"
	useLikes "patreon/internal/app/usecase/likes"
	usePayments "patreon/internal/app/usecase/payments"
	usePosts "patreon/internal/app/usecase/posts"
	usePostsData "patreon/internal/app/usecase/posts_data"
	useSubscr "patreon/internal/app/usecase/subscribers"
	useUser "patreon/internal/app/usecase/user"
)

type UsecaseFactory struct {
	repositoryFactory  RepositoryFactory
	userUsecase        useUser.Usecase
	creatorUsecase     useCreator.Usecase
	csrfUsecase        usecase_csrf.Usecase
	sessionsManager    sessions.SessionsManager
	accessUsecase      useAccess.Usecase
	subscribersUsecase useSubscr.Usecase
	awardsUsercase     useAwards.Usecase
	awardsUsecase      useAwards.Usecase
	sessinManager      sessions.SessionsManager
	postsUsecase       usePosts.Usecase
	postsDataUsecase   usePostsData.Usecase
	likesUsecase       useLikes.Usecase
	paymentsUsecase    usePayments.Usecase
}

func NewUsecaseFactory(repositoryFactory RepositoryFactory) *UsecaseFactory {
	return &UsecaseFactory{
		repositoryFactory: repositoryFactory,
	}
}

func (f *UsecaseFactory) GetUserUsecase() useUser.Usecase {
	if f.userUsecase == nil {
		f.userUsecase = useUser.NewUserUsecase(f.repositoryFactory.GetUserRepository(), f.repositoryFactory.GetFilesRepository())
	}
	return f.userUsecase
}

func (f *UsecaseFactory) GetCreatorUsecase() useCreator.Usecase {
	if f.creatorUsecase == nil {
		f.creatorUsecase = useCreator.NewCreatorUsecase(f.repositoryFactory.GetCreatorRepository(),
			f.repositoryFactory.GetFilesRepository())
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
func (f *UsecaseFactory) GetSubscribersUsecase() useSubscr.Usecase {
	if f.subscribersUsecase == nil {
		f.subscribersUsecase = useSubscr.NewSubscribersUsecase(f.repositoryFactory.GetSubscribersRepository(),
			f.repositoryFactory.GetAwardsRepository())
	}
	return f.subscribersUsecase
}

func (f *UsecaseFactory) GetAwardsUsecase() useAwards.Usecase {
	if f.awardsUsecase == nil {
		f.awardsUsecase = useAwards.NewAwardsUsecase(f.repositoryFactory.GetAwardsRepository(),
			f.repositoryFactory.GetFilesRepository())
	}
	return f.awardsUsecase
}

func (f *UsecaseFactory) GetPostsUsecase() usePosts.Usecase {
	if f.postsUsecase == nil {
		f.postsUsecase = usePosts.NewPostsUsecase(f.repositoryFactory.GetPostsRepository(),
			f.repositoryFactory.GetPostsDataRepository(), f.repositoryFactory.GetFilesRepository())
	}
	return f.postsUsecase
}

func (f *UsecaseFactory) GetLikesUsecase() useLikes.Usecase {
	if f.likesUsecase == nil {
		f.likesUsecase = useLikes.NewLikesUsecase(f.repositoryFactory.GetLikesRepository())
	}
	return f.likesUsecase
}

func (f *UsecaseFactory) GetPostsDataUsecase() usePostsData.Usecase {
	if f.postsDataUsecase == nil {
		f.postsDataUsecase = usePostsData.NewPostsDataUsecase(f.repositoryFactory.GetPostsDataRepository(),
			f.repositoryFactory.GetFilesRepository())
	}
	return f.postsDataUsecase
}
func (f *UsecaseFactory) GetPaymentsUsecase() usePayments.Usecase {
	if f.paymentsUsecase == nil {
		f.paymentsUsecase = usePayments.NewPaymentsUsecase(f.repositoryFactory.GetPaymentsRepository())
	}
	return f.paymentsUsecase
}
