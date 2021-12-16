package repository_subscribers

import "patreon/internal/app/models"

//go:generate mockgen -destination=mocks/mock_subscribers_repository.go -package=mock_repository -mock_names=Repository=SubscribersRepository . Repository

type Repository interface {
	// Create Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	Create(subscriber *models.Subscriber, payToken string) error
	// Delete Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	Delete(subscriber *models.Subscriber) error
	// GetCreators Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	GetCreators(userID int64) ([]models.CreatorSubscribe, error)
	// GetSubscribers Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	GetSubscribers(creatorID int64) ([]models.User, error)
	// Get Errors:
	//		app.GeneralError with Errors
	//			repository.DefaultErrDB
	Get(subscriber *models.Subscriber) (bool, error)
}
