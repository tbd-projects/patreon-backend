package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

type PostsRepository struct {
	store *sql.DB
}

func NewPostsRepository(st *sql.DB) *PostsRepository {
	return &PostsRepository{
		store: st,
	}
}

// Create Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *PostsRepository) Create(post *models.Posts) (int64, error) {
	if err := repo.store.QueryRow("INSERT INTO posts (title, description, "+
		"type_awards, creator_id, cover) VALUES ($1, $2, $3, $4, $5)"+
		"RETURNING creator_id", post.Title, post.Description, post.Awards, post.CreatorId, post.Cover).
		Scan(&post.ID); err != nil {
		return -1, repository.NewDBError(err)
	}
	return post.ID, nil
}

// GetPost Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetPost(postID int64) (*models.Posts, error) {
	post := &models.Posts{ID: postID}

	rows, err := repo.store.Query(
		"SELECT title, description, likes, date, cover, aw.name " +
			"from posts join awards as aw on aw.creator_id = posts.creator_id")
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var creator models.Creator
		if err = rows.Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
			return nil, repository.NewDBError(err)
		}
		res[i] = creator
		i++

		if err = rows.Err(); err != nil {
			return nil, repository.NewDBError(err)
		}
	}
	if err = rows.Close(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// GetPosts Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetPosts(page int64) ([]models.Creator, error) {
	count := 0

	if err := repo.store.QueryRow("SELECT count(*) from creator_profile").Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	res := make([]models.Creator, count)

	rows, err := repo.store.Query(
		"SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
			"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id")
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var creator models.Creator
		if err = rows.Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
			return nil, repository.NewDBError(err)
		}
		res[i] = creator
		i++

		if err = rows.Err(); err != nil {
			return nil, repository.NewDBError(err)
		}
	}
	if err = rows.Close(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// GetCreator Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetCreator(creatorId int64) (*models.Creator, error) {
	creator := &models.Creator{}

	if err := repo.store.QueryRow("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname "+
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id where creator_id=$1", creatorId).
		Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return creator, nil
}
