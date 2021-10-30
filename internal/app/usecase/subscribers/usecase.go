package usecase_subscribers

import "patreon/internal/app/models"

type Usecase interface {
	// Subscribe Errors:
	//		SubscriptionAlreadyExists
	//		app.generalError with Errors
	//			repository.DefaultErrDB
	Subscribe(subscriber *models.Subscriber) error
	// UnSubscribe Errors:
	//		SubscriptionsNotFound
	//		app.generalError with Errors
	//			repository.DefaultErrDB
	UnSubscribe(subscriber *models.Subscriber) error
	// GetCreators Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	GetCreators(userID int64) ([]int64, error)
	// GetSubscribers Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	GetSubscribers(creatorID int64) ([]int64, error)
}
