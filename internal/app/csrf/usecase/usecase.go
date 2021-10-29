package usecase_csrf

import "patreon/internal/app/csrf/models"

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
	Create(sessionId string, userId int64) (models.Token, error)
}
