package handler_factory

import (
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/creator_handler"
	"patreon/internal/app/delivery/http/handlers/creator_handler/subscribe_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler"
	aw_subscribe_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler/subscribe_handler"
	aw_upd_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler/upd_aw_handler"
	upd_cover_awards_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler/upd_cover_awards"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/likes_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/posts_data_id_handler"
	upd_img_data_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/posts_data_id_handler/upd_image_post_handler"
	upd_text_data_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/posts_data_id_handler/upd_text_post_handler"
	upl_cover_posts_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/upd_cover_post_handler"
	posts_upd_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/upd_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/upl_img_data_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/upl_text_data_handler"
	upd_avatar_creator_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/upd_avatar_handler"
	upd_cover_creator_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/upd_cover_handler"
	"patreon/internal/app/delivery/http/handlers/csrf_handler"
	"patreon/internal/app/delivery/http/handlers/info_handler"
	"patreon/internal/app/delivery/http/handlers/login_handler"
	"patreon/internal/app/delivery/http/handlers/logout_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/payments_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/subscriptions_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/avatar_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/nickname_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/password_handler"
	"patreon/internal/app/delivery/http/handlers/register_handler"
	"patreon/internal/microservices/auth/delivery/grpc/client"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
)

const (
	ROOT = iota
	INFO
	REGISTER
	LOGIN
	LOGOUT
	PROFILE
	CREATORS
	CREATOR_WITH_ID
	CREATOR_AVATAR
	CREATOR_COVER
	UPDATE_PASSWORD
	UPDATE_AVATAR
	UPDATE_NICKNAME
	AWARDS
	AWARDS_WITH_ID
	AWARDS_UPDATE
	AWARDS_COVER
	POSTS
	POSTS_WITH_ID
	POSTS_UPD
	POSTS_LIKES
	GET_CSRF_TOKEN
	GET_USER_SUBSCRIPTIONS
	POST_UPD_COVER
	POST_ADD_TEXT
	POST_ADD_IMAGE
	POST_DATA_UPD_TEXT
	POST_DATA_UPD_IMAGE
	POST_DATA_ID
	SUBSCRIBES
	AWARDS_CREATOR_SUBSCRIBE
	USER_PAYMENTS
)

type HandlerFactory struct {
	usecaseFactory    UsecaseFactory
	sessionClientConn *grpc.ClientConn
	logger            *logrus.Logger
	urlHandler        *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, usecaseFactory UsecaseFactory, sClientConn *grpc.ClientConn) *HandlerFactory {
	return &HandlerFactory{
		usecaseFactory:    usecaseFactory,
		logger:            logger,
		sessionClientConn: sClientConn,
	}
}

func (f *HandlerFactory) initAllHandlers() map[int]app.Handler {
	ucUser := f.usecaseFactory.GetUserUsecase()
	ucCreator := f.usecaseFactory.GetCreatorUsecase()
	ucCsrf := f.usecaseFactory.GetCsrfUsecase()
	ucAwards := f.usecaseFactory.GetAwardsUsecase()
	ucPosts := f.usecaseFactory.GetPostsUsecase()
	ucLikes := f.usecaseFactory.GetLikesUsecase()
	ucSubscr := f.usecaseFactory.GetSubscribersUsecase()
	ucPostsData := f.usecaseFactory.GetPostsDataUsecase()
	ucPayments := f.usecaseFactory.GetPaymentsUsecase()
	ucInfo := f.usecaseFactory.GetInfoUsecase()
	sManager := client.NewSessionClient(f.sessionClientConn)

	return map[int]app.Handler{
		INFO:                     info_handler.NewInfoHandler(f.logger, ucInfo),
		REGISTER:                 register_handler.NewRegisterHandler(f.logger, sManager, ucUser),
		LOGIN:                    login_handler.NewLoginHandler(f.logger, sManager, ucUser),
		LOGOUT:                   logout_handler.NewLogoutHandler(f.logger, sManager),
		PROFILE:                  profile_handler.NewProfileHandler(f.logger, sManager, ucUser),
		CREATORS:                 creator_handler.NewCreatorHandler(f.logger, sManager, ucCreator, ucUser),
		CREATOR_WITH_ID:          creator_id_handler.NewCreatorIdHandler(f.logger, sManager, ucCreator),
		UPDATE_PASSWORD:          password_handler.NewUpdatePasswordHandler(f.logger, sManager, ucUser),
		UPDATE_AVATAR:            avatar_handler.NewUpdateAvatarHandler(f.logger, sManager, ucUser),
		UPDATE_NICKNAME:          nickname_handler.NewUpdateNicknameHandler(f.logger, sManager, ucUser),
		AWARDS:                   aw_handler.NewAwardsHandler(f.logger, ucAwards, sManager),
		AWARDS_WITH_ID:           aw_id_handler.NewAwardsIdHandler(f.logger, ucAwards, sManager),
		AWARDS_UPDATE:            aw_upd_handler.NewAwardsUpdHandler(f.logger, ucAwards, sManager),
		POSTS:                    posts_handler.NewPostsHandler(f.logger, ucPosts, sManager),
		POSTS_WITH_ID:            posts_id_handler.NewPostsIDHandler(f.logger, ucPosts, sManager),
		POSTS_UPD:                posts_upd_handler.NewPostsUpdateHandler(f.logger, ucPosts, sManager),
		POSTS_LIKES:              likes_handler.NewLikesHandler(f.logger, ucLikes, ucPosts, sManager),
		GET_CSRF_TOKEN:           csrf_handler.NewCsrfHandler(f.logger, sManager, ucCsrf),
		GET_USER_SUBSCRIPTIONS:   subscriptions_handler.NewSubscriptionsHandler(f.logger, sManager, ucSubscr),
		SUBSCRIBES:               subscribe_handler.NewSubscribeHandler(f.logger, sManager, ucSubscr),
		POST_UPD_COVER:           upl_cover_posts_handler.NewPostsUpdateCoverHandler(f.logger, ucPosts, sManager),
		POST_ADD_TEXT:            upl_text_data_handler.NewPostsDataUploadTextHandler(f.logger, ucPostsData, ucPosts, sManager),
		POST_ADD_IMAGE:           upl_img_data_handler.NewPostsUploadImageHandler(f.logger, ucPostsData, ucPosts, sManager),
		POST_DATA_ID:             posts_data_id_handler.NewPostsDataIDHandler(f.logger, ucPostsData, ucPosts, sManager),
		CREATOR_AVATAR:           upd_avatar_creator_handler.NewUpdateAvatarHandler(f.logger, sManager, ucCreator),
		CREATOR_COVER:            upd_cover_creator_handler.NewUpdateCoverHandler(f.logger, sManager, ucCreator),
		AWARDS_COVER:             upd_cover_awards_handler.NewUpdateCoverAwardsHandler(f.logger, sManager, ucAwards),
		POST_DATA_UPD_IMAGE:      upd_img_data_handler.NewPostsUploadImageHandler(f.logger, ucPostsData, ucPosts, sManager),
		POST_DATA_UPD_TEXT:       upd_text_data_handler.NewPostsDataUpdateTextHandler(f.logger, ucPostsData, ucPosts, sManager),
		AWARDS_CREATOR_SUBSCRIBE: aw_subscribe_handler.NewAwardsSubscribeHandler(f.logger, sManager, ucSubscr, ucAwards),
		USER_PAYMENTS:            payments_handler.NewPaymentsHandler(f.logger, sManager, ucPayments),
	}
}

func (f *HandlerFactory) GetHandleUrls() *map[string]app.Handler {
	if f.urlHandler != nil {
		return f.urlHandler
	}

	hs := f.initAllHandlers()
	f.urlHandler = &map[string]app.Handler{
		//"/":                     "I am a joke?",
		"/login":    hs[LOGIN],
		"/info":     hs[INFO],
		"/logout":   hs[LOGOUT],
		"/register": hs[REGISTER],
		// /user     ---------------------------------------------------------////
		"/user":                 hs[PROFILE],
		"/user/update/password": hs[UPDATE_PASSWORD],
		"/user/update/avatar":   hs[UPDATE_AVATAR],
		"/user/update/nickname": hs[UPDATE_NICKNAME],
		"/user/subscriptions":   hs[GET_USER_SUBSCRIPTIONS],
		"/user/payments":        hs[USER_PAYMENTS],
		// /creators ---------------------------------------------------------////
		"/creators":                                   hs[CREATORS],
		"/creators/{creator_id:[0-9]+}":               hs[CREATOR_WITH_ID],
		"/creators/{creator_id:[0-9]+}/subscribers":   hs[SUBSCRIBES],
		"/creators/{creator_id:[0-9]+}/update/avatar": hs[CREATOR_AVATAR],
		"/creators/{creator_id:[0-9]+}/update/cover":  hs[CREATOR_COVER],
		// ../awards ---------------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/awards":                                hs[AWARDS],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}":              hs[AWARDS_WITH_ID],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}/update":       hs[AWARDS_UPDATE],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}/update/cover": hs[AWARDS_COVER],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}/subscribe":    hs[AWARDS_CREATOR_SUBSCRIBE],
		// ../posts  ---------------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/posts":                         hs[POSTS],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}":        hs[POSTS_WITH_ID],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/update": hs[POSTS_UPD],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/text":   hs[POST_ADD_TEXT],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/image":  hs[POST_ADD_IMAGE],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/cover":  hs[POST_UPD_COVER],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/like":   hs[POSTS_LIKES],
		// ../posts_data  ----------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{data_id:[0-9]+}":              hs[POST_DATA_ID],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{data_id:[0-9]+}/update/text":  hs[POST_DATA_UPD_TEXT],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{data_id:[0-9]+}/update/image": hs[POST_DATA_UPD_IMAGE],
		//   /token  ---------------------------------------------------------////
		"/token": hs[GET_CSRF_TOKEN],
	}
	return f.urlHandler
}
