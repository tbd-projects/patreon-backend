package repository_jwt

import "patreon/internal/app/csrf/models"

//go:generate mockgen -destination=mocks/mock_jwt_repository.go -package=mock_repository -mock_names=Repository=JwtRepository . Repository

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
