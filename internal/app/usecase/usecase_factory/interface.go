package usecase_factory

import (
	repCsrf "patreon/internal/app/csrf/repository/jwt"
	repAccess "patreon/internal/app/repository/access"
	repoAttaches "patreon/internal/app/repository/attaches"
	repoAwrds "patreon/internal/app/repository/awards"
	repoComments "patreon/internal/app/repository/comments"
	repCreator "patreon/internal/app/repository/creator"
	repoInfo "patreon/internal/app/repository/info"
	repoLikes "patreon/internal/app/repository/likes"
	repoPayToken "patreon/internal/app/repository/pay_token"
	repoPayments "patreon/internal/app/repository/payments"
	repoPosts "patreon/internal/app/repository/posts"
	repoStats "patreon/internal/app/repository/statistics"
	useSubscr "patreon/internal/app/repository/subscribers"
	repUser "patreon/internal/app/repository/user"
)

//go:generate mockgen -destination=mocks/mock_repository_factory.go -package=mock_repository_factory . RepositoryFactory

type RepositoryFactory interface {
	GetUserRepository() repUser.Repository
	GetCreatorRepository() repCreator.Repository
	GetAwardsRepository() repoAwrds.Repository
	GetCsrfRepository() repCsrf.Repository
	GetAccessRepository() repAccess.Repository
	GetSubscribersRepository() useSubscr.Repository
	GetPostsRepository() repoPosts.Repository
	GetLikesRepository() repoLikes.Repository
	GetAttachesRepository() repoAttaches.Repository
	GetPaymentsRepository() repoPayments.Repository
	GetInfoRepository() repoInfo.Repository
	GetCommentsRepository() repoComments.Repository
	GetStatsRepository() repoStats.Repository
	GetPayTokenRepository() repoPayToken.Repository
}
