package usecase_csrf

import "patreon/internal/app/csrf/csrf_models"

//go:generate mockgen -destination=mocks/mock_csrf_usecase.go -package=mock_usecase_csrf -mock_names=Usecase=CsrfUsecase . Usecase

type Usecase interface {
	// Check Errors:
	//      repository_jwt.BadToken
	// 		app.GeneralError with Error
	// 			repository_jwt.ParseClaimsError
	//			repository_jwt.TokenExpired
	Check(sessionId string, userId int64, token string) error

	// Create Errors:
	// 		app.GeneralError with Error
	// 			repository_jwt.ErrorSignedToken
	Create(sessionId string, userId int64) (csrf_models.Token, error)
}
