package usecase_subscribers

import "patreon/internal/app/models"

type Usecase interface {
	// Subscribe Errors:
	//		SubscriptionAlreadyExists
	//		repository_postgresql.AwardNameNotFound
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
	GetCreators(userID int64) ([]models.Creator, error)
	// GetSubscribers Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	GetSubscribers(creatorID int64) ([]models.User, error)
}
