package repository

import (
	"errors"
	"time"
)

//go:generate easyjson -disallow_unknown_fields interface.go

var NotModify = errors.New("Not modified push ")

//easyjson:json
type PushJson struct {
	Push interface{} `json:"push"`
}


type PaymentsInfo struct {
	CreatorId  int64
	UserId     int64
	AwardsId   int64
	AwardsName string
}

type Push struct {
	Id     int64
	Type   string
	Push   PushJson
	Date   time.Time
	Viewed bool
}

type Repository interface {
	// GetCreatorNameAndAvatar Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreatorNameAndAvatar(creatorId int64) (nickname string, avatar string, err error)

	// GetUserNameAndAvatar Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetUserNameAndAvatar(userId int64) (nickname string, avatar string, err error)

	// GetAwardsNameAndPrice Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwardsNameAndPrice(awardsId int64) (name string, price int64, err error)

	// GetCreatorPostAndTitle Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreatorPostAndTitle(postId int64) (int64, string, error)

	// GetSubUserForPushPost Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetSubUserForPushPost(postId int64) ([]int64, error)

	// CheckCreatorForGetSubPush Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CheckCreatorForGetSubPush(creatorId int64) (bool, error)

	// GetAwardsInfoAndCreatorIdAndUserIdFromPayments Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwardsInfoAndCreatorIdAndUserIdFromPayments(token string) (*PaymentsInfo, error)

	// CheckCreatorForGetCommentPush Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CheckCreatorForGetCommentPush(creatorId int64) (bool, error)

	// AddPushInfo Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	AddPushInfo(userId []int64, pushType string, push interface{}) error

	// GetPushInfo Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPushInfo(userId int64) ([]Push, error)

	// MarkViewed Errors:
	//		repository.NotModify
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	MarkViewed(pushId int64, userId int64) error
}
