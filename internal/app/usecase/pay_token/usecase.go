package usecase_pay_token

import "patreon/internal/app/models"

//go:generate mockgen -destination=mocks/mock_pay_token_usecase.go -package=mock_usecase -mock_names=Usecase=PayTokenUsecase . Usecase

type Usecase interface {
	// GetToken with Errors:
	//		app.GeneralError with Errors
	//			repository_redis.SetError
	GetToken(userID int64) (models.PayToken, error)
}
