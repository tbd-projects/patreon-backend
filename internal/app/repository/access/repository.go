package repository_access

//go:generate mockgen -destination=mocks/mock_access_repository.go -package=mock_repository -mock_names=Repository=AccessRepository . Repository

type Repository interface {
	// Set Errors:
	// 		app.GeneralError with Errors
	// 			SetError
	Set(key string, value string, timeExp int) error

	// Get Errors:
	//		NotFound
	// 		app.GeneralError with Errors
	// 			InvalidStorageData
	Get(key string) (string, error)

	// Increment Errors:
	// 		app.GeneralError with Errors
	// 			InvalidStorageData
	Increment(userIp string) (int64, error)
}
