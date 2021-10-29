package usecase_factory

import (
	repoAwrds "patreon/internal/app/repository/awards"
	repCreator "patreon/internal/app/repository/creator"
	repoLikes "patreon/internal/app/repository/likes"
	repoPosts "patreon/internal/app/repository/posts"
	repoPostsData "patreon/internal/app/repository/posts_data"
	repUser "patreon/internal/app/repository/user"
	"patreon/internal/app/sessions"
)

type RepositoryFactory interface {
	GetUserRepository() repUser.Repository
	GetCreatorRepository() repCreator.Repository
	GetAwardsRepository() repoAwrds.Repository
	GetSessionRepository() sessions.SessionRepository
	GetPostsRepository() repoPosts.Repository
	GetLikesRepository() repoLikes.Repository
	GetPostsDataRepository() repoPostsData.Repository
}
