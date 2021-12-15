package http_models

import (
	"fmt"
	"math"
	"patreon/internal/app/csrf/csrf_models"
	"patreon/internal/app/models"
	"strconv"
	"time"
)

//go:generate easyjson -all -disallow_unknown_fields response_models.go

//easyjson:json
type TokenResponse struct {
	Token csrf_models.Token `json:"token"`
}

//easyjson:json
type PayTokenResponse struct {
	Token string `json:"token"`
}

//easyjson:json
type PayAccountResponse struct {
	Account string `json:"account_number"`
}

//easyjson:json
type ErrResponse struct {
	Err string `json:"error"`
}

//easyjson:json
type OkResponse struct {
	Ok string `json:"OK"`
}

//easyjson:json
type IdResponse struct {
	ID int64 `json:"id"`
}

//easyjson:json
type ProfileResponse struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	HaveCreator bool   `json:"have_creator"`
}

//easyjson:json
type ResponseInfo struct {
	models.Info
}

//easyjson:json
type ResponseCreatorWithAwards struct {
	models.CreatorWithAwards
}

//easyjson:json
type ResponseCreator struct {
	models.Creator
}

//easyjson:json
type ResponseCreators struct {
	Creators []ResponseCreator `json:"creators"`
}

//easyjson:json
type ResponseCreatorSubscrube struct {
	models.CreatorSubscribe
}

//easyjson:json
type ResponseAward struct {
	ID          int64  `json:"awards_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Color       Color  `json:"color,omitempty"`
	Cover       string `json:"cover"`
	ChildAward  int64  `json:"child_award,omitempty"`
}

//easyjson:json
type ResponseAwards struct {
	Awards []ResponseAward
}

//easyjson:json
type ResponsePost struct {
	ID          int64     `json:"posts_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Awards      int64     `json:"type_awards,omitempty"`
	Likes       int64     `json:"likes"`
	Cover       string    `json:"cover"`
	AddLike     bool      `json:"add_like,omitempty"`
	Views       int64     `json:"views"`
	Comments    int64     `json:"comments"`
	Date        time.Time `json:"date"`
	IsDraft     bool      `json:"is_draft,omitempty"`
}

//easyjson:json
type ResponsePosts struct {
	Posts []ResponsePost
}

//easyjson:json
type ResponsePostComment struct {
	ID             int64     `json:"comment_id"`
	Body           string    `json:"body"`
	AsCreator      bool      `json:"as_creator,omitempty"`
	AuthorId       int64     `json:"author_id"`
	Date           time.Time `json:"date"`
	AuthorNickname string    `json:"author_nickname"`
	AuthorAvatar   string    `json:"author_avatar"`
}

//easyjson:json
type ResponseUserComment struct {
	ID        int64     `json:"comment_id"`
	Body      string    `json:"body"`
	AsCreator bool      `json:"as_creator,omitempty"`
	PostId    int64     `json:"post_id"`
	Date      time.Time `json:"date"`
	PostName  string    `json:"post_name"`
	PostCover string    `json:"post_cover"`
}

//easyjson:json
type ResponseUserComments struct {
	Comments []ResponseUserComment `json:"comments"`
}

//easyjson:json
type ResponsePostComments struct {
	Comments []ResponsePostComment `json:"comments"`
}

//easyjson:json
type ResponseAttach struct {
	ID    int64  `json:"attach_id"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

//easyjson:json
type ResponseApplyAttach struct {
	IDs []int64 `json:"attaches_id"`
}

//easyjson:json
type ResponsePostWithAttaches struct {
	Post ResponsePost     `json:"post"`
	Data []ResponseAttach `json:"attaches"`
}

//easyjson:json
type ResponseBalance struct {
	ID      int64        `json:"user_id"`
	Balance models.Money `json:"balance"`
}

func ToRProfileResponse(us models.User) ProfileResponse {
	return ProfileResponse{
		ID:          us.ID,
		Login:       us.Login,
		Nickname:    us.Nickname,
		Avatar:      us.Avatar,
		HaveCreator: us.HaveCreator,
	}
}

func ToResponseUserComment(cm models.UserComment) ResponseUserComment {
	return ResponseUserComment{
		ID:        cm.ID,
		Body:      cm.Body,
		PostId:    cm.PostId,
		Date:      cm.Date,
		PostName:  cm.PostName,
		PostCover: cm.PostCover,
		AsCreator: cm.AsCreator,
	}
}

func ToResponsePostComment(cm models.PostComment) ResponsePostComment {
	return ResponsePostComment{
		ID:             cm.ID,
		Body:           cm.Body,
		AuthorId:       cm.AuthorId,
		Date:           cm.Date,
		AuthorNickname: cm.AuthorNickname,
		AuthorAvatar:   cm.AuthorAvatar,
		AsCreator:      cm.AsCreator,
	}
}

func ToResponsePostComments(cms []models.PostComment) ResponsePostComments {
	res := ResponsePostComments{[]ResponsePostComment{}}
	for _, cm := range cms {
		res.Comments = append(res.Comments, ToResponsePostComment(cm))
	}
	return res
}

func ToResponseUserComments(cms []models.UserComment) ResponseUserComments {
	res := ResponseUserComments{[]ResponseUserComment{}}
	for _, cm := range cms {
		res.Comments = append(res.Comments, ToResponseUserComment(cm))
	}
	return res
}

func ToResponseCreator(cr models.Creator) ResponseCreator {
	return ResponseCreator{
		cr,
	}
}

func ToResponseInfo(info models.Info) ResponseInfo {
	return ResponseInfo{
		info,
	}
}

func ToResponseCreators(crs []models.Creator) ResponseCreators {
	respondCreators := make([]ResponseCreator, len(crs))
	for i, cr := range crs {
		respondCreators[i] = ToResponseCreator(cr)
	}
	return ResponseCreators{respondCreators}
}

func ToResponseCreatorWithAwards(cr models.CreatorWithAwards) ResponseCreatorWithAwards {
	res := ResponseCreatorWithAwards{
		cr,
	}
	res.AwardsId = int64(math.Max(float64(res.AwardsId), 0))
	return res
}

func ToResponseAward(aw models.Award) ResponseAward {
	return ResponseAward{
		ID:          aw.ID,
		Name:        aw.Name,
		Price:       aw.Price,
		Description: aw.Description,
		Color:       NewColor(aw.Color),
		Cover:       aw.Cover,
		ChildAward:  int64(math.Max(float64(aw.ChildAward), 0)),
	}
}

func ToResponsePost(ps models.Post) ResponsePost {
	return ResponsePost{
		ID:          ps.ID,
		Title:       ps.Title,
		Description: ps.Description,
		Date:        ps.Date,
		Likes:       ps.Likes,
		Awards:      int64(math.Max(float64(ps.Awards), 0)),
		Cover:       ps.Cover,
		AddLike:     ps.AddLike,
		Views:       ps.Views,
		IsDraft:     ps.IsDraft,
		Comments:    ps.Comments,
	}
}

func ToResponsePostWithAttaches(ps models.PostWithAttach) ResponsePostWithAttaches {
	res := ResponsePostWithAttaches{Post: ToResponsePost(*ps.Post), Data: []ResponseAttach{}}
	for _, data := range ps.Data {
		res.Data = append(res.Data, ToResponseAttach(data))
	}
	return res
}

func ToResponseAttach(ps models.AttachWithoutLevel) ResponseAttach {
	return ResponseAttach{
		ID:    ps.ID,
		Value: ps.Value,
		Type:  string(ps.Type),
	}
}

func (u *ResponseCreator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s}", strconv.Itoa(int(u.ID)), u.Nickname)
}

//easyjson:json
type SubscriptionsUserResponse struct {
	Creators []ResponseCreatorSubscrube `json:"creators"`
}

func ToSubscriptionsUser(creators []models.CreatorSubscribe) SubscriptionsUserResponse {
	var res []ResponseCreatorSubscrube
	for _, creator := range creators {
		res = append(res, ResponseCreatorSubscrube{
			models.CreatorSubscribe{
				ID:          creator.ID,
				Nickname:    creator.Nickname,
				Description: creator.Description,
				Category:    creator.Category,
				Cover:       creator.Cover,
				Avatar:      creator.Avatar,
				AwardsId:    creator.AwardsId,
			},
		})
	}
	return SubscriptionsUserResponse{
		Creators: res,
	}
}

//easyjson:json
type ResponseUser struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

//easyjson:json
type SubscribersCreatorResponse struct {
	Users []ResponseUser `json:"users"`
}

func ToSubscribersCreatorResponse(users []models.User) SubscribersCreatorResponse {
	res := make([]ResponseUser, len(users))
	for i, user := range users {
		res[i] = ResponseUser{
			ID:       user.ID,
			Login:    user.Login,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		}
	}
	return SubscribersCreatorResponse{
		Users: res,
	}
}

//easyjson:json
type ResponseLike struct {
	Likes int64 `json:"likes"`
}

//easyjson:json
type ResponseUserPayments struct {
	Payments []models.UserPayments `json:"payments"`
}
type ResponseCreatorPayments struct {
	Payments []models.CreatorPayments `json:"payments"`
}

func ToResponseUserPayments(payments []models.UserPayments) ResponseUserPayments {
	res := make([]models.UserPayments, 0, len(payments))
	for _, payment := range payments {
		res = append(res, models.UserPayments{
			Payments: models.Payments{
				Amount:    payment.Amount,
				Date:      payment.Date,
				CreatorID: payment.CreatorID,
				Status:    payment.Status,
			},
			CreatorNickname:    payment.CreatorNickname,
			CreatorDescription: payment.CreatorDescription,
			CreatorCategory:    payment.CreatorCategory,
		})
	}
	return ResponseUserPayments{
		Payments: res,
	}
}
func ToResponseCreatorPayments(payments []models.CreatorPayments) ResponseCreatorPayments {
	res := make([]models.CreatorPayments, 0, len(payments))
	for _, payment := range payments {
		res = append(res, models.CreatorPayments{
			Payments: models.Payments{
				Amount: payment.Amount,
				Date:   payment.Date,
				UserID: payment.UserID,
				Status: payment.Status,
			},
			UserNickname: payment.UserNickname,
		})
	}
	return ResponseCreatorPayments{
		Payments: res,
	}
}

//easyjson:json
type ResponseAvailablePosts struct {
	AvailablePosts []models.AvailablePost `json:"available_posts"`
}

func ToResponseAvailablePosts(availablePosts []models.AvailablePost) ResponseAvailablePosts {
	return ResponseAvailablePosts{
		AvailablePosts: availablePosts,
	}
}

//easyjson:json
type ResponseCreatorPostsViews struct {
	CountPostsViews int64 `json:"count_posts_views"`
}

//easyjson:json
type ResponseCreatorCountSubscribers struct {
	CountSubscribers int64 `json:"count_subscribers"`
}

//easyjson:json
type ResponseCreatorTotalIncome struct {
	TotalIncome float64 `json:"total_income"`
}

//easyjson:json
type ResponseCreatorCountPosts struct {
	CountPosts int64 `json:"count_posts"`
}

//easyjson:json
type ResponsePayToken struct {
	PayToken string `json:"token"`
}

//easyjson:json
type ResponsePayAccount struct {
	Account string `json:"account_number"`
}
