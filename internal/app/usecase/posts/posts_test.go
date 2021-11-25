package posts

import "patreon/internal/app/usecase"

type SuitePostsUsecase struct {
	usecase.SuiteUsecase
	uc Usecase
}
