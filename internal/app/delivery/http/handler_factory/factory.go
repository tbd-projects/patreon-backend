package handler_factory

import (
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/creator_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/awards_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/awards_handler/awards_id_handler"
	aw_other_update_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/awards_handler/awards_id_handler/update_handler/update_awards_other_handler"
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
	GET_CSRF_TOKEN
	GET_USER_SUBSCRIPTIONS
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
		AWARDS:                 awards_handler.NewAwardsHandler(f.logger, f.router, f.cors, ucAwards, sManager),
		AWARDS_WITH_ID:         awards_id_handler.NewAwardsIDHandler(f.logger, f.router, f.cors, ucAwards, sManager),
		AWARDS_OTHER_UPD:       aw_other_update_handler.NewAwardsUpOtherHandler(f.logger, f.router, f.cors, ucAwards, sManager),
		GET_CSRF_TOKEN:         csrf_handler.NewCsrfHandler(f.logger, f.router, f.cors, sManager, ucCsrf),
		GET_USER_SUBSCRIPTIONS: subscriptions_handler.NewSubscriptionsHandler(f.logger, f.router, f.cors, sManager, ucSubscr),
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
		"/creators":                     hs[CREATORS],
		"/creators/{creator_id:[0-9]+}": hs[CREATOR_WITH_ID],
		// ../awards ---------------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/awards":                                 hs[AWARDS],
		"/creators/{creator_id:[0-9]+}/awards/{awards_id:[0-9]+}":              hs[AWARDS_WITH_ID],
		"/creators/{creator_id:[0-9]+}/awards/{awards_id:[0-9]+}/update/other": hs[AWARDS_OTHER_UPD],
		// ../posts  ---------------------------------------------------------////
		//   /token  ---------------------------------------------------------////
		"/token": hs[GET_CSRF_TOKEN],
	}
	return f.urlHandler
}
