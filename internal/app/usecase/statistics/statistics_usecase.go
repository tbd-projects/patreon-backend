package statistics

import (
	"patreon/internal/app"
	repository_statistics "patreon/internal/app/repository/statistics"
)

type StatisticsUsecase struct {
	repository repository_statistics.Repository
}

func NewStatisticsUsecase(repository repository_statistics.Repository) *StatisticsUsecase {
	return &StatisticsUsecase{
		repository: repository,
	}
}

// GetCountCreatorPosts Errors:
//		CreatorDoesNotExists
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (u *StatisticsUsecase) GetCountCreatorPosts(creatorID int64) (int64, error) {
	isExists, err := u.repository.CreatorExists(creatorID)
	if err != nil {
		return app.InvalidInt, err
	}

	if !isExists {
		return app.InvalidInt, CreatorDoesNotExists
	}

	return u.repository.GetCountCreatorPosts(creatorID)
}

// GetCountCreatorSubscribers Errors:
//		CreatorDoesNotExists
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (u *StatisticsUsecase) GetCountCreatorSubscribers(creatorID int64) (int64, error) {
	isExists, err := u.repository.CreatorExists(creatorID)
	if err != nil {
		return app.InvalidInt, err
	}

	if !isExists {
		return app.InvalidInt, CreatorDoesNotExists
	}

	return u.repository.GetCountCreatorSubscribers(creatorID)
}

// GetCountCreatorViews Errors:
//		CreatorDoesNotExists
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (u *StatisticsUsecase) GetCountCreatorViews(creatorID int64, days int64) (int64, error) {
	isExists, err := u.repository.CreatorExists(creatorID)
	if err != nil {
		return app.InvalidInt, err
	}

	if !isExists {
		return app.InvalidInt, CreatorDoesNotExists
	}

	return u.repository.GetCountCreatorViews(creatorID, days)
}

// GetTotalIncome Errors:
//		CreatorDoesNotExists
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (u *StatisticsUsecase) GetTotalIncome(creatorID int64, days int64) (float64, error) {
	isExists, err := u.repository.CreatorExists(creatorID)
	if err != nil {
		return app.InvalidFloat, err
	}

	if !isExists {
		return app.InvalidFloat, CreatorDoesNotExists
	}

	return u.repository.GetTotalIncome(creatorID, days)
}
