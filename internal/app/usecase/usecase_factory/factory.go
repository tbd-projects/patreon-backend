package usecase_factory

import (
	usecase_csrf "patreon/internal/app/csrf/usecase"
	useAccess "patreon/internal/app/usecase/access"
	useAttaches "patreon/internal/app/usecase/attaches"
	useAwards "patreon/internal/app/usecase/awards"
	useComments "patreon/internal/app/usecase/comments"
	useCreator "patreon/internal/app/usecase/creator"
	useInfo "patreon/internal/app/usecase/info"
	useLikes "patreon/internal/app/usecase/likes"
	usePayToken "patreon/internal/app/usecase/pay_token"
	usePayments "patreon/internal/app/usecase/payments"
	usePosts "patreon/internal/app/usecase/posts"
	useStats "patreon/internal/app/usecase/statistics"
	useSubscr "patreon/internal/app/usecase/subscribers"
	useUser "patreon/internal/app/usecase/user"
	"patreon/internal/microservices/files/delivery/grpc/client"

	"google.golang.org/grpc"
)

type UsecaseFactory struct {
	repositoryFactory  RepositoryFactory
	userUsecase        useUser.Usecase
	creatorUsecase     useCreator.Usecase
	csrfUsecase        usecase_csrf.Usecase
	accessUsecase      useAccess.Usecase
	subscribersUsecase useSubscr.Usecase
	awardsUsecase      useAwards.Usecase
	postsUsecase       usePosts.Usecase
	attachesUsecase    useAttaches.Usecase
	infoUsecase        useInfo.Usecase
	likesUsecase       useLikes.Usecase
	paymentsUsecase    usePayments.Usecase
	fileClient         client.FileServiceClient
	statsUsecase       useStats.Usecase
	commentsUsecase    useComments.Usecase
	payTokenUsecase    usePayToken.Usecase
}

func NewUsecaseFactory(repositoryFactory RepositoryFactory, fileConn *grpc.ClientConn) *UsecaseFactory {
	fileClient := client.NewFileServiceClient(fileConn)
	return &UsecaseFactory{
		repositoryFactory: repositoryFactory,
		fileClient:        fileClient,
	}
}

func (f *UsecaseFactory) GetUserUsecase() useUser.Usecase {
	if f.userUsecase == nil {
		f.userUsecase = useUser.NewUserUsecase(f.repositoryFactory.GetUserRepository(), f.fileClient)
	}
	return f.userUsecase
}

func (f *UsecaseFactory) GetCreatorUsecase() useCreator.Usecase {
	if f.creatorUsecase == nil {
		f.creatorUsecase = useCreator.NewCreatorUsecase(f.repositoryFactory.GetCreatorRepository(),
			f.fileClient)
	}
	return f.creatorUsecase
}
func (f *UsecaseFactory) GetCsrfUsecase() usecase_csrf.Usecase {
	if f.csrfUsecase == nil {
		f.csrfUsecase = usecase_csrf.NewCsrfUsecase(f.repositoryFactory.GetCsrfRepository())
	}
	return f.csrfUsecase
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
			f.fileClient)
	}
	return f.awardsUsecase
}

func (f *UsecaseFactory) GetPostsUsecase() usePosts.Usecase {
	if f.postsUsecase == nil {
		f.postsUsecase = usePosts.NewPostsUsecase(f.repositoryFactory.GetPostsRepository(),
			f.repositoryFactory.GetAttachesRepository(), f.fileClient)
	}
	return f.postsUsecase
}

func (f *UsecaseFactory) GetLikesUsecase() useLikes.Usecase {
	if f.likesUsecase == nil {
		f.likesUsecase = useLikes.NewLikesUsecase(f.repositoryFactory.GetLikesRepository())
	}
	return f.likesUsecase
}

func (f *UsecaseFactory) GetAttachesUsecase() useAttaches.Usecase {
	if f.attachesUsecase == nil {
		f.attachesUsecase = useAttaches.NewAttachesUsecase(f.repositoryFactory.GetAttachesRepository(),
			f.fileClient)
	}
	return f.attachesUsecase
}

func (f *UsecaseFactory) GetPaymentsUsecase() usePayments.Usecase {
	if f.paymentsUsecase == nil {
		f.paymentsUsecase = usePayments.NewPaymentsUsecase(f.repositoryFactory.GetPaymentsRepository())
	}
	return f.paymentsUsecase
}

func (f *UsecaseFactory) GetInfoUsecase() useInfo.Usecase {
	if f.infoUsecase == nil {
		f.infoUsecase = useInfo.NewInfoUsecase(f.repositoryFactory.GetInfoRepository())
	}
	return f.infoUsecase
}
func (f *UsecaseFactory) GetStatsUsecase() useStats.Usecase {
	if f.statsUsecase == nil {
		f.statsUsecase = useStats.NewStatisticsUsecase(f.repositoryFactory.GetStatsRepository())
	}
	return f.statsUsecase
}

func (f *UsecaseFactory) GetCommentsUsecase() useComments.Usecase {
	if f.commentsUsecase == nil {
		f.commentsUsecase = useComments.NewCommentsUsecase(f.repositoryFactory.GetCommentsRepository())
	}
	return f.commentsUsecase
}
func (f *UsecaseFactory) GetPayTokenUsecase() usePayToken.Usecase {
	if f.payTokenUsecase == nil {
		f.payTokenUsecase = usePayToken.NewPayTokenUsecase(f.repositoryFactory.GetPayTokenRepository())
	}
	return f.payTokenUsecase
}
