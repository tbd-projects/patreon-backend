package usecase_access

import (
	"patreon/internal/app/models"
	repository_access "patreon/internal/app/repository/access"
	"time"
)

var (
	queryLimit  int64 = 100
	timeLimit   int   = int(time.Minute.Milliseconds())
	blackList         = "BLACK_LIST"
	timeBlocked       = int(time.Minute.Milliseconds() * 2)
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

	accessCounter := models.AccessCounter{
		Counter: userCounter,
		Limit:   queryLimit,
	}

	if accessCounter.Overflow() {
		return false, NoAccess
	}

	return true, nil
}

// Create Errors:
// 		app.GeneralError with Errors
// 			repository_access.SetError
func (u *AccessUsecase) Create(userIp string) (bool, error) {
	if err := u.repository.Set(userIp, *models.NewAccessCounter(queryLimit), timeLimit); err != nil {
		return false, err
	}
	return true, nil
}

// Update Errors:
// 		app.GeneralError with Errors
// 			repository_access.InvalidStorageData
func (u *AccessUsecase) Update(userIp string) (bool, error) {
	if ok, err := u.repository.Update(userIp); err != nil || ok == -1 {
		return false, err
	}
	return true, nil
}

// AddToBlackList Errors:
// 		app.GeneralError with Errors
// 			repository_access.SetError
func (u *AccessUsecase) AddToBlackList(userIp string) error {
	if err := u.repository.AddToBlackList(blackList, userIp, timeBlocked); err != nil {
		return err
	}
	return nil
}

// CheckBlackList Errors:
// 		app.GeneralError with Errors
// 			repository_access.InvalidStorageData
func (u *AccessUsecase) CheckBlackList(userIp string) (bool, error) {
	_, err := u.repository.CheckBlackList(blackList, userIp)
	switch err {
	case repository_access.NotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return true, err
	}
}
