package repository_postgresql

import (
	"database/sql"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_posts "patreon/internal/app/repository/likes"
	rp "patreon/internal/app/repository/posts"
	putilits "patreon/internal/app/utilits/postgresql"
)

type PostsRepository struct {
	store *sql.DB
}

var _ = repository_posts.Repository(&PostsRepository{})

func NewPostsRepository(st *sql.DB) *PostsRepository {
	return &PostsRepository{
		store: st,
	}
}

// Create Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *PostsRepository) Create(post *models.Post) (int64, error) {
	query := `INSERT INTO posts (title, description,
		type_awards, creator_id, cover) VALUES ($1, $2, $3, $4, $5) 
		RETURNING posts_id`
	var awardsId putilits.NullInt64
	awardsId.Int64 = post.Awards
	if post.Awards == rp.NoAwards {
		awardsId.Valid = false
	} else {
		awardsId.Valid = true
	}

	if err := repo.store.QueryRow(query, post.Title, post.Description, awardsId, post.CreatorId, "").
		Scan(&post.ID); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}
	return post.ID, nil
}

// GetPost Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetPost(postID int64) (*models.Post, error) {
	query := `SELECT title, description, likes, date, cover, awards_id FROM posts WHERE post_id = $1`
	post := &models.Post{ID: postID}
	var awardsId putilits.NullInt64
	if err := repo.store.QueryRow(query, postID).Scan(&post.Title, &post.Description,
		&post.Likes, &post.Date, &post.Date, &awardsId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	if awardsId.Valid == false {
		post.Awards = rp.NoAwards
	} else {
		post.Awards = awardsId.Int64
	}

	return post, nil
}

// GetPosts Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetPosts(creatorsId int64, pag models.Pagination) ([]models.Creator, error) {
	queryCount := `SELECT count(*) FROM posts WHERE posts`
	queryPost := `SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname 
					FROM creator_profile JOIN users AS usr ON usr.user_id = creator_profile.creator_id`
	queryPost =
	count := 0

	if err := repo.store.QueryRow(queryCount).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	res := make([]models.Creator, count)

	rows, err := repo.store.Query(queryPost)
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

// UpdatePost Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) UpdatePost(post *models.Post) error {
	query := `UPDATE posts SET title = $1, description = $2, type_awards = $3 WHERE posts_id = $4`

	var awardsId putilits.NullInt64
	awardsId.Int64 = post.Awards
	if post.Awards == rp.NoAwards {
		awardsId.Valid = false
	} else {
		awardsId.Valid = true
	}

	if _, err := repo.store.Query(query, post.Title, post.Description, awardsId, post.ID); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// UpdateCoverPost Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) UpdateCoverPost(postId int64, cover string) error {
	query := `UPDATE posts SET cover = $1 WHERE posts_id = $2`

	if _, err := repo.store.Query(query, cover, postId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) Delete(postId int64) error {
	query := `DELETE FROM posts WHERE posts_id = $q`

	if _, err := repo.store.Query(query, postId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
