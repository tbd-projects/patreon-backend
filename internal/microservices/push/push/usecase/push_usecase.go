package usecase

import (
	"patreon/internal/microservices/push"
	"patreon/internal/microservices/push/push"
	"patreon/internal/microservices/push/push/repository"
)

type PushUsecase struct {
	repository repository.Repository
}

func NewPushUsecase(repository repository.Repository) *PushUsecase {
	return &PushUsecase{
		repository: repository,
	}
}

// PreparePostPush with Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) PreparePostPush(info push.PostInfo) ([]int64, *push_models.PostPush, error) {
	result := &push_models.PostPush{
		PostId:    info.PostId,
		PostTitle: info.PostTitle,
		CreatorId: info.CreatorId,
	}

	nickname, avatar, err := usecase.repository.GetCreatorNameAndAvatar(info.CreatorId)
	if err != nil {
		return nil, nil, err
	}

	result.CreatorNickname = nickname
	result.CreatorAvatar = avatar

	allow, err := usecase.repository.GetSubUserForPushPost(info.PostId)
	return allow, result, err
}

// PrepareCommentPush with Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) PrepareCommentPush(info push.CommentInfo) ([]int64, *push_models.CommentPush, error) {
	result := &push_models.CommentPush{
		CommentId: info.CommentId,
		AuthorId:  info.AuthorId,
		PostId:    info.PostId,
	}

	nickname, avatar, err := usecase.repository.GetCreatorNameAndAvatar(info.AuthorId)
	if err != nil {
		return nil, nil, err
	}

	result.AuthorNickname = nickname
	result.AuthorAvatar = avatar

	creatroId, title, err := usecase.repository.GetCreatorPostAndTitle(info.PostId)
	if err != nil {
		return nil, nil, err
	}

	result.PostTitle = title

	allow, err := usecase.repository.CheckCreatorForGetCommentPush(creatroId)
	if err != nil {
		return nil, nil, err
	}

	if allow {
		return []int64{creatroId}, result, err
	}
	return []int64{}, result, err
}

// PrepareSubPush with Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) PrepareSubPush(info push.SubInfo) ([]int64, *push_models.SubPush, error) {
	result := &push_models.SubPush{
		UserId:   info.UserId,
		AwardsId: info.AwardsId,
	}

	nickname, avatar, err := usecase.repository.GetCreatorNameAndAvatar(info.UserId)
	if err != nil {
		return nil, nil, err
	}

	result.UserAvatar = nickname
	result.UserAvatar = avatar

	name, price, err := usecase.repository.GetAwardsNameAndPrice(info.AwardsId)
	if err != nil {
		return nil, nil, err
	}

	result.AwardsName = name
	result.AwardsPrice = price

	allow, err := usecase.repository.CheckCreatorForGetCommentPush(info.CreatorId)
	if err != nil {
		return nil, nil, err
	}

	if allow {
		return []int64{info.CreatorId}, result, err
	}
	return []int64{}, result, err
}
