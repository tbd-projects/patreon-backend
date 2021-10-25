package repository_token

import "patreon/internal/app/csrf/models"

type Repository interface {
	// Check Errors:
	// 		repository_jwt.BadToken
	// 		app.GeneralError with Error
	// 			repository_jwt.ParseClaimsError
	// 			repository_jwt.TokenExpired
	Check(sources models.TokenSources, tokenString models.Token) error

	// Create Errors:
	// 		app.GeneralError with Error
	// 			repository_jwt.ErrorSignedToken
	Create(sources models.TokenSources) (models.Token, error)
}
