package handler_factory

import (
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/creator_handler"
	"patreon/internal/app/delivery/http/handlers/creator_handler/subscribe_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_handler/aw_id_handler"
	aw_upd_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_handler/aw_id_handler/upd_aw_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_handler/posts_id_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_handler/posts_id_handler/likes_handler"
	posts_upd_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_handler/posts_id_handler/update_handler"
	"patreon/internal/app/delivery/http/handlers/csrf_handler"
	"patreon/internal/app/delivery/http/handlers/login_handler"
	"patreon/internal/app/delivery/http/handlers/logout_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/subscriptions_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/avatar_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/password_handler"
	handlers2 "patreon/internal/app/delivery/http/handlers/register_handler"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ROOT = iota
	REGISTER
	LOGIN
	LOGOUT
	PROFILE
	CREATORS
	CREATOR_WITH_ID
	UPDATE_PASSWORD
	UPDATE_AVATAR
	AWARDS
	AWARDS_WITH_ID
	AWARDS_OTHER_UPD
	POSTS
	POSTS_WITH_ID
	POSTS_UPD
	POSTS_LIKES
	GET_CSRF_TOKEN
	GET_USER_SUBSCRIPTIONS
	SUBSCRIBES
)

type HandlerFactory struct {
	usecaseFactory UsecaseFactory
	logger         *logrus.Logger
	router         *mux.Router
	cors           *app.CorsConfig
	urlHandler     *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, router *mux.Router,
	cors *app.CorsConfig, usecaseFactory UsecaseFactory) *HandlerFactory {
	return &HandlerFactory{
		usecaseFactory: usecaseFactory,
		logger:         logger,
		router:         router,
		cors:           cors,
	}
}

func (f *HandlerFactory) initAllHandlers() map[int]app.Handler {
	ucUser := f.usecaseFactory.GetUserUsecase()
	ucCreator := f.usecaseFactory.GetCreatorUsecase()
	ucCsrf := f.usecaseFactory.GetCsrfUsecase()
	ucAwards := f.usecaseFactory.GetAwardsUsecase()
	ucPosts := f.usecaseFactory.GetPostsUsecase()
	ucLikes := f.usecaseFactory.GetLikesUsecase()
	sManager := f.usecaseFactory.GetSessionManager()
	ucSubscr := f.usecaseFactory.GetSubscribersUsecase()

	return map[int]app.Handler{
		REGISTER:               handlers2.NewRegisterHandler(f.logger, f.router, f.cors, sManager, ucUser),
		LOGIN:                  login_handler.NewLoginHandler(f.logger, f.router, f.cors, sManager, ucUser),
		LOGOUT:                 logout_handler.NewLogoutHandler(f.logger, f.router, f.cors, sManager),
		PROFILE:                profile_handler.NewProfileHandler(f.logger, f.router, f.cors, sManager, ucUser),
		CREATORS:               creator_handler.NewCreatorHandler(f.logger, f.router, f.cors, sManager, ucCreator, ucUser),
		CREATOR_WITH_ID:        creator_id_handler.NewCreatorIdHandler(f.logger, f.router, f.cors, sManager, ucUser, ucCreator),
		UPDATE_PASSWORD:        password_handler.NewUpdatePasswordHandler(f.logger, f.router, f.cors, sManager, ucUser),
		UPDATE_AVATAR:          avatar_handler.NewUpdateAvatarHandler(f.logger, f.router, f.cors, sManager, ucUser),
		AWARDS:                 aw_handler.NewAwardsHandler(f.logger, f.router, f.cors, ucAwards, sManager),
		AWARDS_WITH_ID:         aw_id_handler.NewAwardsIdHandler(f.logger, f.router, f.cors, ucAwards, sManager),
		AWARDS_OTHER_UPD:       aw_upd_handler.NewAwardsUpdHandler(f.logger, f.router, f.cors, ucAwards, sManager),
		POSTS:                  posts_handler.NewPostsHandler(f.logger, f.router, f.cors, ucPosts, sManager),
		POSTS_WITH_ID:          posts_id_handler.NewPostsIDHandler(f.logger, f.router, f.cors, ucPosts, sManager),
		POSTS_UPD:              posts_upd_handler.NewPostsUpdateHandler(f.logger, f.router, f.cors, ucPosts, sManager),
		POSTS_LIKES:            likes_handler.NewLikesHandler(f.logger, f.router, f.cors, ucLikes, ucPosts, sManager),
		GET_CSRF_TOKEN:         csrf_handler.NewCsrfHandler(f.logger, f.router, f.cors, sManager, ucCsrf),
		GET_USER_SUBSCRIPTIONS: subscriptions_handler.NewSubscriptionsHandler(f.logger, f.router, f.cors, sManager, ucSubscr),
		SUBSCRIBES:             subscribe_handler.NewSubscribeHandler(f.logger, f.router, f.cors, sManager, ucSubscr),
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
		"/logout":   hs[LOGOUT],
		"/register": hs[REGISTER],
		// /user     ---------------------------------------------------------////
		"/user":                 hs[PROFILE],
		"/user/update/password": hs[UPDATE_PASSWORD],
		"/user/update/avatar":   hs[UPDATE_AVATAR],
		"/user/subscriptions":   hs[GET_USER_SUBSCRIPTIONS],
		// /creators ---------------------------------------------------------////
		"/creators":                               hs[CREATORS],
		"/creators/{creator_id:[0-9]+}":           hs[CREATOR_WITH_ID],
		"/creators/{creator_id:[0-9]+}/subscribe": hs[SUBSCRIBES],
		// ../awards ---------------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/awards":                                hs[AWARDS],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}":              hs[AWARDS_WITH_ID],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}/update/other": hs[AWARDS_OTHER_UPD],
		// ../posts  ---------------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/posts":                         hs[POSTS],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}":        hs[POSTS_WITH_ID],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/update": hs[POSTS_UPD],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/like":   hs[POSTS_LIKES],
		//   /token  ---------------------------------------------------------////
		"/token": hs[GET_CSRF_TOKEN],
	}
	return f.urlHandler
}
