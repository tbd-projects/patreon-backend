package repository_access

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

	// Update Errors:
	// 		app.GeneralError with Errors
	// 			InvalidStorageData
	Update(userIp string) (string, error)
}
