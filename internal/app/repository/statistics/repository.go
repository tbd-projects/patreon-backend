package repository_statistics

//go:generate mockgen -destination=mocks/mock_statistics_repository.go -package=mock_repository -mock_names=Repository=StatisticsRepository . Repository

type Repository interface {
	// GetCountCreatorPosts Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetCountCreatorPosts(creatorID int64) (int64, error)
	// GetCountCreatorSubscribers Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetCountCreatorSubscribers(creatorID int64) (int64, error)
	// GetCountCreatorViews Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetCountCreatorViews(creatorID int64, days int64) (int64, error)
	// GetTotalIncome Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetTotalIncome(creatorID int64, days int64) (float64, error)
}
