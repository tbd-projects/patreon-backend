package repository_jwt

import "patreon/internal/app/csrf/models"

type Repository interface {
	// Check Errors:
	// 		BadToken
	// 		app.GeneralError with Error
	// 			ParseClaimsError
	// 			TokenExpired
	Check(sources models.TokenSources, tokenString models.Token) error

	// Create Errors:
	// 		app.GeneralError with Error
	// 			ErrorSignedToken
	Create(sources models.TokenSources) (models.Token, error)
}
