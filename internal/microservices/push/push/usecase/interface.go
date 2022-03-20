package usecase

import (
	"patreon/internal/microservices/push"
	"patreon/internal/microservices/push/push"
	"patreon/internal/microservices/push/push/repository"
)

type Usecase interface {
	// PreparePostPush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PreparePostPush(info *push.PostInfo) ([]int64, *push_models.PostPush, error)

	// PrepareCommentPush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PrepareCommentPush(info *push.CommentInfo) ([]int64, *push_models.CommentPush, error)

	// PreparePaymentsPush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PreparePaymentsPush(info *push.PaymentApply) ([]int64, *push_models.PaymentApplyPush, error)

	// AddPushInfo Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	AddPushInfo(userId []int64, pushType string, push interface{}) error

	// GetPushInfo Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPushInfo(userId int64) ([]repository.Push, error)

	// MarkViewed Errors:
	//		repository.NotModify
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	MarkViewed(pushId int64, userId int64) error
}
