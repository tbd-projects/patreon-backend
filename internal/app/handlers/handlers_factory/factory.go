package handlers_factory

import (
	"fmt"
	"patreon/internal/app"
	"patreon/internal/app/handlers/urls"
)

type Joiner interface {
	JoinHandlers(...app.Joinable)
}

type JoinerAndJoinable interface {
	app.Joinable
	Joiner
}

type HandlersFactory struct {
	useHandlers map[urls.UrlPath]JoinerAndJoinable
	routingUrls map[urls.UrlPath]bool
}

func NewHandlerFactory(useHandlers map[urls.UrlPath]JoinerAndJoinable, routingUrls []urls.UrlPath) HandlersFactory {
	setRoutingUrls := map[urls.UrlPath]bool{}
	for _, url := range routingUrls {
		setRoutingUrls[url] = true
	}
	return HandlersFactory{useHandlers: useHandlers, routingUrls: setRoutingUrls}
}

func (factory *HandlersFactory) AddUseHandler(handlerUrl urls.UrlPath, handler JoinerAndJoinable) {
	factory.useHandlers[handlerUrl] = handler
}

func (factory *HandlersFactory) AddRoutingUrl(url urls.UrlPath) {
	factory.routingUrls[url] = true
}

func (factory *HandlersFactory) JoinHandlers(startHandler Joiner) error {
	var useUrls []urls.UrlPath
	for url, _ := range factory.routingUrls {
		useUrls = append(useUrls, url)
	}

	var allowedUrl []urls.UrlPath
	for url, _ := range factory.routingUrls {
		allowedUrl = append(allowedUrl, url)
	}

	return factory.recursiveJoin(startHandler, useUrls, allowedUrl)
}

func (factory *HandlersFactory) recursiveJoin(currentHandler Joiner,
	nextUrls []urls.UrlPath, allowedUrl []urls.UrlPath) error {

	var nextHandlers map[urls.UrlPath]JoinerAndJoinable
	for _, url := range nextUrls {
		handlerUrl := getFirstUrl(allowedUrl, url)
		nextHandler, ok := factory.useHandlers[handlerUrl]
		if !ok {
			return fmt.Errorf("In routing url there is containg unknown handler url %s ", handlerUrl)
		}
		nextHandlers[handlerUrl] = nextHandler
	}

	for url, handler := range nextHandlers {
		currentHandler.JoinHandlers(handler)
		err := factory.recursiveJoin(handler, getUrlsWithBeginning(nextUrls, url), allowedUrl)
		if err != nil {
			return err
		}
	}
	return nil
}
