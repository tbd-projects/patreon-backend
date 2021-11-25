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
	att := RequestAttaches{Attaches: []RequestAttach{{Value: "sdsd"}}}
	att.Sanitize(*bluemonday.UGCPolicy())
	changeNick := RequestChangeNickname{}
	changeNick.Sanitize(*bluemonday.UGCPolicy())

	_ = ToResponsePost(models.Post{})
	_ = ToResponseAttach(models.AttachWithoutLevel{})
	_ = ToResponseAward(models.Award{})
	_ = ToResponseCreator(models.Creator{})
	_ = ToResponsePostWithAttaches(models.PostWithAttach{Post: &models.Post{}, Data: []models.AttachWithoutLevel{{Value: "sdsd"}}})
	_ = ToRProfileResponse(models.User{})
	_ = ToSubscribersCreatorResponse([]models.User{{Nickname: ""}})
	_ = ToSubscriptionsUser([]models.CreatorSubscribe{{ID: 3}})
	_ = ToResponseAttach(models.AttachWithoutLevel{})
	_ = ToResponseUserPayments([]models.UserPayments{{Payments: models.Payments{Amount: 0}}})
	_ = ToResponseInfo(models.Info{})
	_ = ToResponseCreators([]models.Creator{{Nickname: ""}})
}
