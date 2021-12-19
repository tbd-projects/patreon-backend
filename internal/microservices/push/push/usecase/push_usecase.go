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
func (usecase *PushUsecase) PreparePostPush(info *push.PostInfo) ([]int64, *push_models.PostPush, error) {
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
func (usecase *PushUsecase) PrepareCommentPush(info *push.CommentInfo) ([]int64, *push_models.CommentPush, error) {
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
		result.CreatorId = creatroId
		return []int64{creatroId}, result, err
	}
	return []int64{}, result, err
}

// PreparePaymentsPush with Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) PreparePaymentsPush(info *push.PaymentApply) ([]int64, *push_models.PaymentApplyPush, error) {
	result := &push_models.PaymentApplyPush{}

	payment, err := usecase.repository.GetAwardsInfoAndCreatorIdAndUserIdFromPayments(info.Token)
	if err != nil {
		return nil, nil, err
	}

	result.AwardsId = payment.AwardsId
	result.AwardsName = payment.AwardsName

	nickname, avatar, err := usecase.repository.GetCreatorNameAndAvatar(result.CreatorId)
	if err != nil {
		return nil, nil, err
	}

	result.CreatorNickname = nickname
	result.CreatorAvatar = avatar
	return []int64{payment.UserId}, result, err
}

// AddPushInfo Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) AddPushInfo(userId []int64, pushType string, push interface{}) error {
	return usecase.repository.AddPushInfo(userId, pushType, push)
}

// GetPushInfo Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) GetPushInfo(userId int64) ([]repository.Push, error) {
	return usecase.repository.GetPushInfo(userId)
}

// MarkViewed Errors:
//		repository.NotModify
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) MarkViewed(pushId int64, userId int64) error {
	return usecase.repository.MarkViewed(pushId, userId)
}
