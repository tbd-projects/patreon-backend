package http_models

import (
	"fmt"
	"patreon/internal/app/models"

	"github.com/pkg/errors"
)

var (
	AwardNameValidateError = errors.New("invalid award_name")
	NicknameValidateError  = errors.New(fmt.Sprintf("invalid nickname in body len must be from %v to %v",
		models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH))
)
