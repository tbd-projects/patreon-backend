package repository_user

import (
	"errors"
	"patreon/internal/app/repository"

	"github.com/lib/pq"
)

const (
	codeDuplicateVal   = "23505"
	loginConstraint    = "users_login_key"
	nicknameConstraint = "users_nickname_key"
)

var (
	LoginAlreadyExist    = errors.New("login already exist")
	NicknameAlreadyExist = errors.New("nickname already exist")
)

func parsePQError(err *pq.Error) error {
	switch {
	case err.Code == codeDuplicateVal && err.Constraint == loginConstraint:
		return LoginAlreadyExist
	case err.Code == codeDuplicateVal && err.Constraint == nicknameConstraint:
		return NicknameAlreadyExist
	default:
		return repository.NewDBError(err)
	}
}
