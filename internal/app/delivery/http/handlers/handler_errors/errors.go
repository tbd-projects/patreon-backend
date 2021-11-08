package handler_errors

import (
	"errors"
	"fmt"
	"patreon/internal/app/models"
)

/// NOT FOUND
var (
	UserNotFound             = errors.New("user not found")
	UserWithNicknameNotFound = errors.New("user with this nickname not found")
	AwardNotFound            = errors.New("award with this id not found")
	PostNotFound             = errors.New("post with not found")
	PostDataNotFound         = errors.New("post data with this id not found")
	LikeNotFound             = errors.New("like with this id not found")
	PaymentsNotFound         = errors.New("this user have not payment")
)

/// Fields Incorrect
var (
	InvalidNickname          = errors.New("invalid creator nickname")
	InvalidCategory          = errors.New("invalid creator category")
	InvalidDescription       = errors.New("invalid creator description")
	IncorrectAwardsId        = errors.New("this awards id not know")
	IncorrectPostId          = errors.New("this post id not know")
	IncorrectCreatorId       = errors.New("this creator id not know")
	EmptyTitle               = errors.New("empty title")
	EmptyName                = errors.New("empty name in request")
	IncorrectLoginOrPassword = errors.New("incorrect login or password")
	IncorrectPrice           = errors.New("incorrect value of price")
	IncorrectNewPassword     = errors.New("invalid new password")
	IncorrectDataType        = errors.New("invalid data type")
	InvalidOldNickname       = errors.New("old nickname not equal current user nickname")
)

// BD Error
var (
	LikesAlreadyDel          = errors.New("this user not have like for this post")
	LikesAlreadyExists       = errors.New("this user already add like for this post")
	AwardsAlreadyExists      = errors.New("awards with this name already exists")
	AwardsPriceAlreadyExists = errors.New("awards with this price already exists")
	UserAlreadyExist         = errors.New("user already exist")
	NicknameAlreadyExist     = errors.New("nickname already exist")
	CreatorAlreadyExist      = errors.New("creator already exist")
	BDError                  = errors.New("can not do bd operation")
)

// Session Error
var (
	ErrorCreateSession = errors.New("can not create session")
	DeleteCookieFail   = errors.New("can not delete cookie from session store")
)

// Request Error
var (
	InvalidBody          = errors.New("invalid body in request")
	InvalidParameters    = errors.New("invalid parameters")
	UserNotHaveAward     = errors.New("this user not have award for this post")
	InvalidQueries       = errors.New("invalid parameters in query")
	FileSizeError        = errors.New("size of file very big")
	InvalidFormFieldName = errors.New("invalid form field name for load file")
	InvalidImageExt      = errors.New("please upload a JPEG, JPG or PNG files")
	UserAlreadySubscribe = errors.New("this user already have subscribe on creator")
	SubscribesNotFound   = errors.New("subscribes on the creator not found")
	InvalidUserNickname  = errors.New(fmt.Sprintf("invalid nickname in body len must be from %v to %v",
		models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH))
)

var InternalError = errors.New("server error")
