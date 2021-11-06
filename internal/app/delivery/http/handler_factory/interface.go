package handler_factory

import (
	useCsrf "patreon/internal/app/csrf/usecase"
	"patreon/internal/app/sessions"
	useAwards "patreon/internal/app/usecase/awards"
	useCreator "patreon/internal/app/usecase/creator"
	useLikes "patreon/internal/app/usecase/likes"
	usePayments "patreon/internal/app/usecase/payments"
	usePosts "patreon/internal/app/usecase/posts"
	usePostsData "patreon/internal/app/usecase/posts_data"
	useSubscr "patreon/internal/app/usecase/subscribers"
	useUser "patreon/internal/app/usecase/user"
)

//go:generate mockgen -destination=mocks/mock_usecase_factory.go -package=mock_usecase_factory . UsecaseFactory

type UsecaseFactory interface {
	GetUserUsecase() useUser.Usecase
	GetCreatorUsecase() useCreator.Usecase
	GetCsrfUsecase() useCsrf.Usecase
	GetAwardsUsecase() useAwards.Usecase
	GetPostsUsecase() usePosts.Usecase
	GetSessionManager() sessions.SessionsManager
	GetSubscribersUsecase() useSubscr.Usecase
	GetLikesUsecase() useLikes.Usecase
	GetPostsDataUsecase() usePostsData.Usecase
	GetPaymentsUsecase() usePayments.Usecase
}
