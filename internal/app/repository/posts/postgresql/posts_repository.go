package repository_postgresql

import (
	"database/sql"
	"fmt"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	rp "patreon/internal/app/repository"
	repository_posts "patreon/internal/app/repository/posts"
	putilits "patreon/internal/app/utilits/postgresql"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

const (
	getAvailablePosts = `SELECT p.posts_id, p.title,
			p.description,
			p.likes,
			p.date,
			p.cover,
			p.type_awards,
			p.creator_id,
			u.nickname,
			lk.likes_id IS NOT NULL,
			views,
			p.number_comments
	FROM posts p
			 JOIN subscribers s on s.users_id = $1 and s.creator_id = p.creator_id
			 LEFT JOIN likes AS lk ON (lk.post_id = p.posts_id and lk.users_id = $1)
			 JOIN users u on p.creator_id = u.users_id
	WHERE p.is_draft = false and (p.type_awards is null OR p.type_awards = s.awards_id or p.type_awards in
			 (select awa.awards_id from restapi_dev.public.parents_awards as awa where awa.parent_id = s.awards_id))
	ORDER BY p.date desc LIMIT $2 OFFSET $3;`

	createQuery = `INSERT INTO posts (title, description,
		type_awards, creator_id, cover, is_draft) VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING posts_id`

	getPostCreatorQuery = `SELECT creator_id FROM posts WHERE posts_id = $1`

	getPostQuery = `
			SELECT title, description, likes, posts.date, cover, type_awards, 
			       creator_id, lk.likes_id IS NOT NULL, views, is_draft, number_comments FROM posts
				LEFT OUTER JOIN likes AS lk ON (lk.post_id = posts.posts_id and lk.users_id = $1)
				WHERE posts.posts_id = $2;`
	getPostQueryUpdate = `UPDATE posts SET views = views + 1 WHERE posts_id = $1`

	updateQuery = `UPDATE posts SET title = $1, description = $2, type_awards = $3, is_draft = $4
					WHERE posts_id = $5 RETURNING posts_id`

	updateCoverQuery = `UPDATE posts SET cover = $1 WHERE posts_id = $2 RETURNING posts_id`

	deleteQuery = `DELETE FROM posts WHERE posts_id = $1`

	getPostsQueryWithDraft = `
			SELECT posts_id, title, description, likes, type_awards, posts.date, cover, 
					lk.likes_id IS NOT NULL, views, is_draft, number_comments
			FROM posts
			LEFT JOIN likes AS lk ON (lk.post_id = posts.posts_id and lk.users_id = $1)
			WHERE creator_id = $2 ORDER BY posts.date DESC
	`
	getPostsQueryWithoutDraft = `
			SELECT posts_id, title, description, likes, type_awards, posts.date, cover, 
					lk.likes_id IS NOT NULL, views, number_comments
			FROM posts
			LEFT JOIN likes AS lk ON (lk.post_id = posts.posts_id and lk.users_id = $1)
			WHERE creator_id = $2 AND NOT is_draft ORDER BY posts.date DESC
	`
)

type PostsRepository struct {
	store *sqlx.DB
}

var _ = repository_posts.Repository(&PostsRepository{})

func NewPostsRepository(st *sqlx.DB) *PostsRepository {
	return &PostsRepository{
		store: st,
	}
}

// Create Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *PostsRepository) Create(post *models.CreatePost) (int64, error) {
	var awardsId sql.NullInt64
	awardsId.Int64 = post.Awards
	if post.Awards == rp.NoAwards {
		awardsId.Valid = false
	} else {
		awardsId.Valid = true
	}

	if err := repo.store.QueryRowx(createQuery, post.Title, post.Description, awardsId, post.CreatorId,
		app.DefaultImage, post.IsDraft).
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
	creatorId := int64(0)
	if err := repo.store.QueryRowx(getPostCreatorQuery, postID).Scan(&creatorId); err != nil {
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
	post := &models.Post{ID: postID}
	var awardsId sql.NullInt64
	if err := repo.store.QueryRow(getPostQuery, userId, postID).Scan(&post.Title, &post.Description,
		&post.Likes, &post.Date, &post.Cover, &awardsId,
		&post.CreatorId, &post.AddLike, &post.Views, &post.IsDraft, &post.Comments); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	if addView {
		row, err := repo.store.Query(getPostQueryUpdate, postID)
		if err != nil {
			return nil, repository.NewDBError(err)
		}
		if err = row.Close(); err != nil {
			return nil, repository.NewDBError(err)
		}
	}

	if !awardsId.Valid {
		post.Awards = rp.NoAwards
	} else {
		post.Awards = awardsId.Int64
	}

	return post, nil
}

// GetAvailablePosts Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetAvailablePosts(userID int64, pag *models.Pagination) ([]models.AvailablePost, error) {
	limit, offset, err := putilits.AddPagination("posts", pag, repo.store)

	if err != nil {
		return nil, err
	}

	var res []models.AvailablePost

	rows, err := repo.store.Query(getAvailablePosts, userID, limit, offset)
	if err != nil {
		return nil, repository.NewDBError(err)
	}
	for rows.Next() {
		var post models.AvailablePost
		var awardsId sql.NullInt64
		err = rows.Scan(
			&post.ID, &post.Title, &post.Description, &post.Likes, &post.Date,
			&post.Cover, &awardsId, &post.CreatorId, &post.CreatorNickname,
			&post.AddLike, &post.Views, &post.Comments)

		if err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		if !awardsId.Valid {
			post.Awards = rp.NoAwards
		} else {
			post.Awards = awardsId.Int64
		}

		res = append(res, post)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// GetPosts Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) GetPosts(creatorsId int64, userId int64,
	pag *models.Pagination, withDraft bool) ([]models.Post, error) {

	query := getPostsQueryWithoutDraft
	if withDraft {
		query = getPostsQueryWithDraft
	}
	limit, offset, err := putilits.AddPagination("posts", pag, repo.store)
	query = query + fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	if err != nil {
		return nil, err
	}
	res := make([]models.Post, 0, limit)

	rows, err := repo.store.Query(query, userId, creatorsId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for rows.Next() {
		var post models.Post
		var awardsId sql.NullInt64

		if withDraft {
			err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.Likes,
				&awardsId, &post.Date, &post.Cover, &post.AddLike, &post.Views, &post.IsDraft, &post.Comments)
		} else {
			err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.Likes,
				&awardsId, &post.Date, &post.Cover, &post.AddLike, &post.Views, &post.Comments)
		}

		if err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		if !awardsId.Valid {
			post.Awards = rp.NoAwards
		} else {
			post.Awards = awardsId.Int64
		}
		post.CreatorId = creatorsId

		res = append(res, post)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// UpdatePost Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsRepository) UpdatePost(post *models.UpdatePost) error {
	var awardsId sql.NullInt64
	awardsId.Int64 = post.Awards
	if post.Awards == rp.NoAwards {
		awardsId.Valid = false
	} else {
		awardsId.Valid = true
	}

	var postsId int64
	if err := repo.store.QueryRow(updateQuery, post.Title, post.Description,
		awardsId, post.IsDraft, post.ID).Scan(&postsId); err != nil {
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
	var postsId int64
	if err := repo.store.QueryRow(updateCoverQuery, cover, postId).Scan(&postsId); err != nil {
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
	row, err := repo.store.Query(deleteQuery, postId)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
