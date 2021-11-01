package repository_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_posts "patreon/internal/app/repository/posts"
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
func (repo *PostsRepository) Create(post *models.CreatePost) (int64, error) {
	query := `INSERT INTO posts (title, description,
		type_awards, creator_id, cover) VALUES ($1, $2, $3, $4, $5) 
		RETURNING posts_id`
	var awardsId sql.NullInt64
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

// GetPostCreator Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetPostCreator(postID int64) (int64, error) {
	query := `SELECT creator_id FROM posts WHERE posts_id = $1`
	creatorId := int64(0)
	if err := repo.store.QueryRow(query, postID).Scan(&creatorId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, repository.NotFound
		}
		return app.InvalidInt, repository.NewDBError(err)
	}

	return creatorId, nil
}

// GetPost Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetPost(postID int64, userId int64, addView bool) (*models.Post, error) {
	query := `
			SELECT title, description, likes, posts.date, cover, type_awards, lk.likes_id IS NOT NULL, views FROM posts
				LEFT OUTER JOIN likes AS lk ON (lk.post_id = posts.posts_id and lk.users_id = $1)
				WHERE posts.posts_id = $2;`
	queryPost := `UPDATE posts SET views = views + 1 WHERE posts_id = $1`

	post := &models.Post{ID: postID}
	var awardsId sql.NullInt64
	if err := repo.store.QueryRow(query, userId, postID).Scan(&post.Title, &post.Description,
		&post.Likes, &post.Date, &post.Cover, &awardsId, &post.AddLike, &post.Views); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	if addView {
		row, err := repo.store.Query(queryPost, postID)
		if err != nil {
			return nil, repository.NewDBError(err)
		}
		if err = row.Close(); err != nil {
			return nil, repository.NewDBError(err)
		}
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
func (repo *PostsRepository) GetPosts(creatorsId int64, userId int64, pag *models.Pagination) ([]models.Post, error) {
	limit, offset, err := putilits.AddPagination("posts", pag, repo.store)
	query := `
			SELECT posts_id, title, description, likes, type_awards, posts.date, cover, lk.likes_id IS NOT NULL, views
			FROM posts
			LEFT JOIN likes AS lk ON (lk.post_id = posts.posts_id and lk.users_id = $1)
			WHERE creator_id = $2 ORDER BY posts.date` + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	if err != nil {
		return nil, err
	}
	var res []models.Post

	rows, err := repo.store.Query(query, userId, creatorsId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for rows.Next() {
		var post models.Post
		var awardsId sql.NullInt64
		if err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.Likes,
			&awardsId, &post.Date, &post.Cover, &post.AddLike, &post.Views); err != nil {
			return nil, repository.NewDBError(err)
		}
		if awardsId.Valid == false {
			post.Awards = rp.NoAwards
		} else {
			post.Awards = awardsId.Int64
		}
		post.CreatorId = creatorsId
		res = append(res, post)

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
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) UpdatePost(post *models.UpdatePost) error {
	query := `UPDATE posts SET title = $1, description = $2, type_awards = $3 WHERE posts_id = $4 RETURNING posts_id`

	var awardsId sql.NullInt64
	awardsId.Int64 = post.Awards
	if post.Awards == rp.NoAwards {
		awardsId.Valid = false
	} else {
		awardsId.Valid = true
	}

	var postsId int64
	if err := repo.store.QueryRow(query, post.Title, post.Description, awardsId, post.ID).Scan(&postsId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// UpdateCoverPost Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) UpdateCoverPost(postId int64, cover string) error {
	query := `UPDATE posts SET cover = $1 WHERE posts_id = $2 RETURNING posts_id`

	var postsId int64
	if err := repo.store.QueryRow(query, cover, postId).Scan(&postsId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) Delete(postId int64) error {
	query := `DELETE FROM posts WHERE posts_id = $q`

	row, err := repo.store.Query(query, postId)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
