package repository_user

import (
	"errors"
	"patreon/internal/app/repository"

	"github.com/lib/pq"
)

var (
	LoginAlreadyExist    = errors.New("login already exist")
	NicknameAlreadyExist = errors.New("nickname already exist")
)

func parseDBError(err *pq.Error) error {
	switch {
	case err.Code == "23505" && err.Constraint == "users_login_key":
		return LoginAlreadyExist
	case err.Code == "23505" && err.Constraint == "users_nickname_key":
		return NicknameAlreadyExist
	default:
		return repository.NewDBError(err)
	}
}
