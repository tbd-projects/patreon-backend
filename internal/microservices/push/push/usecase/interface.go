package usecase

import (
	"patreon/internal/microservices/push"
	"patreon/internal/microservices/push/push"
)

type Usecase interface {
	// PreparePostPush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PreparePostPush(info push.PostInfo) ([]int64, *push_models.PostPush, error)

	// PrepareCommentPush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PrepareCommentPush(info push.CommentInfo) ([]int64, *push_models.CommentPush, error)

	// PrepareSubPush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PrepareSubPush(info push.SubInfo) ([]int64, *push_models.SubPush, error)
}
