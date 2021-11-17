package sessions

import "patreon/internal/app/sessions/models"

//go:generate mockgen -destination=mocks/manager_mock.go -package=mock_sessions . SessionsManager
//go:generate mockgen -destination=mocks/repository_mock.go -package=mock_sessions . SessionRepository

type SessionRepository interface {
	Set(session *models.Session) error
	GetUserId(key string) (string, error)
	Del(session *models.Session) error
}

type SessionsManager interface {
	Check(uniqID string) (models.Result, error)
	Create(userID int64) (models.Result, error)
	Delete(uniqID string) error
}
