package http_models

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/require"
	"patreon/internal/app/models"
	"testing"
)

func TestSanitize(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	creator := RequestCreator{}
	creator.Sanitize(*bluemonday.UGCPolicy())
	login := RequestLogin{}
	login.Sanitize(*bluemonday.UGCPolicy())
	pas := RequestChangePassword{}
	pas.Sanitize(*bluemonday.UGCPolicy())
	reg := RequestRegistration{}
	reg.Sanitize(*bluemonday.UGCPolicy())
	awards := RequestAwards{}
	awards.Sanitize(*bluemonday.UGCPolicy())
	pots := RequestPosts{}
	pots.Sanitize(*bluemonday.UGCPolicy())
	text := RequestText{}
	text.Sanitize(*bluemonday.UGCPolicy())
	sub := SubscribeRequest{}
	sub.Sanitize(*bluemonday.UGCPolicy())

	_ = ToResponsePost(models.Post{})
	_ = ToResponseAttach(models.PostData{})
	_ = ToResponseAward(models.Award{})
	_ = ToResponseCreator(models.Creator{})
	_ = ToResponsePostWithAttaches(models.PostWithData{Post:&models.Post{}})
	_ = ToRProfileResponse(models.User{})
	_ = ToSubscribersCreatorResponse([]models.User{})
	_ = ToSubscriptionsUser([]models.CreatorSubscribe{})

}
