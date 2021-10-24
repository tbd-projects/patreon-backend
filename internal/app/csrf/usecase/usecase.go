package usecase_csrf

import "patreon/internal/app/csrf/models"

type Usecase interface {
	Check(sessionId string, userId int64, token string) error
	Create(sessionId string, userId int64) (models.Token, error)
}
