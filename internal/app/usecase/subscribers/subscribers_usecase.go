package usecase_subscribers

import (
	"patreon/internal/app/models"
	repository_subscribers "patreon/internal/app/repository/subscribers"

	"github.com/pkg/errors"
)

type SubscribersUsecase struct {
	repository repository_subscribers.Repository
}

func NewSubscribersUsecase(repository repository_subscribers.Repository) *SubscribersUsecase {
	return &SubscribersUsecase{
		repository: repository,
	}
}

// Create Errors:
//		SubscriptionAlreadyExists
//		app.generalError with Errors
//			repository.DefaultErrDB
func (uc *SubscribersUsecase) Create(subscriber *models.Subscriber) error {
	exist, err := uc.repository.Get(subscriber.UserID, subscriber.CreatorID)
	if err != nil {
		return errors.Wrapf(err, "METHOD: subscribers_usecase.Create; "+
			"ERR: error on checkExists userID = %v creatorID = %v", subscriber.UserID, subscriber.CreatorID)
	}
	if exist {
		return SubscriptionAlreadyExists
	}
	return uc.repository.Create(subscriber)
}

// GetCreators Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (uc *SubscribersUsecase) GetCreators(userID int64) ([]int64, error) {
	return uc.repository.GetCreators(userID)
}

// GetSubscribers Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (uc *SubscribersUsecase) GetSubscribers(creatorID int64) ([]int64, error) {
	return uc.repository.GetSubscribers(creatorID)
}
