package repository_token

import "patreon/internal/app/csrf/models"

type Repository interface {
	Check(sources models.TokenSources, tokenString models.Token) error
	Create(sources models.TokenSources) (models.Token, error)
}
