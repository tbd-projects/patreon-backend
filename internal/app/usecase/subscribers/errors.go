package usecase_subscribers

import "errors"

var (
	SubscriptionAlreadyExists = errors.New("this subscribe already exists")
	SubscriptionsNotFound     = errors.New("the user is not subscribed on this creator")
)
