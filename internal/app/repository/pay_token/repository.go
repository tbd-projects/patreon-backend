package pay_token

//go:generate mockgen -destination=mocks/mock_pay_token_repository.go -package=mock_repository -mock_names=Repository=PayTokenRepository . Repository

type Repository interface {
	// Set Errors:
	// 		app.GeneralError with Errors
	// 			repository_redis.SetError
	Set(key string, value string, timeExp int) error
	// Get Errors:
	//		repository_redis.NotFound
	// 		app.GeneralError with Errors
	// 			repository_redis.InvalidStorageData
	Get(key string) (string, error)
}
