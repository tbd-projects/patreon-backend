package statistics

//go:generate mockgen -destination=mocks/mock_statistics_usecase.go -package=mock_usecase -mock_names=Usecase=StatisticsUsecase . Usecase

type Usecase interface {
	// GetCountCreatorPosts Errors:
	//		CreatorDoesNotExists
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetCountCreatorPosts(creatorID int64) (int64, error)

	// GetCountCreatorSubscribers Errors:
	//		CreatorDoesNotExists
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetCountCreatorSubscribers(creatorID int64) (int64, error)

	// GetCountCreatorViews Errors:
	//		CreatorDoesNotExists
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetCountCreatorViews(creatorID int64, days int64) (int64, error)

	// GetTotalIncome Errors:
	//		CreatorDoesNotExists
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetTotalIncome(creatorID int64, days int64) (float64, error)
}
