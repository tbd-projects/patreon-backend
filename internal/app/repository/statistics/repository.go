package repository_statistics

//go:generate mockgen -destination=mocks/mock_statistics_repository.go -package=mock_repository -mock_names=Repository=StatisticsRepository . Repository

type Repository interface {
	GetCountCreatorPosts(creatorID int64) (int64, error)
	GetCountCreatorSubscribers(creatorID int64) (int64, error)
	GetCountCreatorViews(creatorID int64, days int64) (int64, error)
	GetTotalIncome(creatorID int64, days int64) (int64, error)
}
