package usecase_csrf

import (
	"patreon/internal/app/csrf/models"
	repository_token "patreon/internal/app/csrf/repository"
	"time"
)

var expiredCsrfTime = time.Minute * 15

type CsrfUsecase struct {
	repository repository_token.Repository
}

func NewCsrfUsecase(repo repository_token.Repository) *CsrfUsecase {
	return &CsrfUsecase{
		repository: repo,
	}
}

func (u *CsrfUsecase) Check(sessionId string, userId int64, token models.Token) error {
	sources := models.TokenSources{
		UserId:    userId,
		SessionId: sessionId,
	}
	return u.repository.Check(sources, token)

}
func (u *CsrfUsecase) Create(sessionId string, userId int64) (models.Token, error) {
	data := models.TokenSources{
		UserId:      userId,
		SessionId:   sessionId,
		ExpiredTime: time.Now().Add(expiredCsrfTime),
	}
	return u.repository.Create(data)
}
