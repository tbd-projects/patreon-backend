package usecase_access

type Usecase interface {
	// CheckAccess Errors:
	//		NoAccess
	//		FirstQuery
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	CheckAccess(userIp string) (bool, error)

	// Create Errors:
	// 		app.GeneralError with Errors
	// 			repository_access.SetError
	Create(userIp string) (bool, error)

	// Update Errors:
	//		strconv.NumError
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	Update(userIp string) (int, error)

	// AddToBlackList Errors:
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	AddToBlackList(userIp string) error

	// CheckBlackList Errors:
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	CheckBlackList(userIp string) (bool, error)
}
