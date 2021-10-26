package repository_access

import "patreon/internal/app/models"

type Repository interface {
	// Set Errors:
	// 		app.GeneralError with Errors
	// 			SetError
	Set(userIp string, access models.AccessCounter, timeExp int) error

	// Get Errors:
	//		NotFound
	// 		app.GeneralError with Errors
	// 			InvalidStorageData
	Get(userIp string) (int64, error)

	// Update Errors:
	// 		app.GeneralError with Errors
	// 			InvalidStorageData
	Update(userIp string) (int64, error)

	// AddToBlackList Errors:
	// 		app.GeneralError with Errors
	// 			SetError
	AddToBlackList(key string, userIp string, timeLimit int) error

	// CheckBlackList Errors:
	//		NotFound
	// 		app.GeneralError with Errors
	// 			InvalidStorageData
	CheckBlackList(key string, userIp string) (bool, error)
}
