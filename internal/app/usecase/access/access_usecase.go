package usecase_access

import (
	"patreon/internal/app"
	repository_access "patreon/internal/app/repository/access"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	queryLimit      = 2
	timeLimit   int = int(time.Minute.Milliseconds())
	blackList       = "BLACK_LIST"
	timeBlocked     = int(time.Minute.Milliseconds() * 2)
)

type AccessUsecase struct {
	repository repository_access.Repository
}

func NewAccessUsecase(repo repository_access.Repository) *AccessUsecase {
	return &AccessUsecase{
		repository: repo,
	}
}

// CheckAccess Errors:
//      NoAccess
//		FirstQuery
// 		app.GeneralError with Errors
// 			repository_access.InvalidStorageData
func (u *AccessUsecase) CheckAccess(userIp string) (bool, error) {
	userCounter, err := u.repository.Get(userIp)
	if err == repository_access.NotFound {
		return true, FirstQuery
	}

	if err != nil {
		return false, err
	}
	accessCounter, err := strconv.Atoi(userCounter)
	if err != nil {
		return false, app.GeneralError{
			Err: repository_access.InvalidStorageData,
			ExternalErr: errors.Wrapf(err, "AccessUsecase: can not convert string from repo to int string: %v",
				accessCounter),
		}
	}
	if accessCounter >= queryLimit {
		return false, NoAccess
	}

	return true, nil
}

// Create Errors:
// 		app.GeneralError with Errors
// 			repository_access.SetError
func (u *AccessUsecase) Create(userIp string) (bool, error) {
	countAccesses := "0"
	if err := u.repository.Set(userIp, countAccesses, timeLimit); err != nil {
		return false, err
	}
	return true, nil
}

// Update Errors:
//		NoAccess
// 		app.GeneralError with Errors
// 			repository_access.InvalidStorageData
func (u *AccessUsecase) Update(userIp string) (int64, error) {
	num, err := u.repository.Increment(userIp)
	if err != nil {
		return -1, err
	}
	if int(num) >= queryLimit {
		return -1, NoAccess
	}
	return num, nil
}

// AddToBlackList Errors:
// 		app.GeneralError with Errors
// 			repository_access.SetError
func (u *AccessUsecase) AddToBlackList(userIp string) error {
	blackListKey := blackList + userIp
	if err := u.repository.Set(blackListKey, userIp, timeBlocked); err != nil {
		return err
	}
	return nil
}

// CheckBlackList Errors:
// 		app.GeneralError with Errors
// 			repository_access.InvalidStorageData
func (u *AccessUsecase) CheckBlackList(userIp string) (bool, error) {
	blackListKey := blackList + userIp
	_, err := u.repository.Get(blackListKey)

	switch err {
	case repository_access.NotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return true, err
	}
}
