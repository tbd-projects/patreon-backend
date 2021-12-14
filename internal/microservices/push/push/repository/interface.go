package repository

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

	// CheckCreatorForGetCommentPush Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CheckCreatorForGetCommentPush(creatorId int64) (bool, error)
}
