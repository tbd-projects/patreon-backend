package repository_postgresql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"

	"github.com/pkg/errors"
)

type LikesRepository struct {
	store *sqlx.DB
}

func NewLikesRepository(st *sqlx.DB) *LikesRepository {
	return &LikesRepository{
		store: st,
	}
}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *LikesRepository) Get(userId int64) (*models.Like, error) {
	query := `SELECT post_id, value, likes_id FROM likes WHERE users_id = $1`
	like := &models.Like{UserId: userId}
	if err := repo.store.QueryRow(query, userId).Scan(&like.PostId, &like.Value,
		&like.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return like, nil
}

// GetLikeId Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *LikesRepository) GetLikeId(userId int64, postId int64) (int64, error) {
	query := `SELECT likes_id FROM likes WHERE users_id = $1 AND post_id = $2`
	likeId := int64(0)
	if err := repo.store.QueryRow(query, userId, postId).Scan(&likeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, repository.NotFound
		}
		return app.InvalidInt, repository.NewDBError(err)
	}

	return likeId, nil
}

// Add Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *LikesRepository) Add(like *models.Like) (int64, error) {
	queryInsert := `INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0)`
	queryUpdate := `UPDATE posts SET likes = likes + $2 WHERE posts_id = $1 RETURNING likes;`

	begin, err := repo.store.Begin()
	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}
	var row *sql.Rows
	row, err = begin.Query(queryInsert, like.PostId, like.UserId, like.Value)

	if err != nil {
		_ = begin.Rollback()
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return app.InvalidInt, repository.NewDBError(err)
	}
	countLikes := int64(app.InvalidInt)
	err = begin.QueryRow(queryUpdate, like.PostId, like.Value).Scan(&countLikes)

	if err != nil {
		_ = begin.Rollback()
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err = begin.Commit(); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}
	return countLikes, nil
}

// Delete Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *LikesRepository) Delete(likeId int64) (int64, error) {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`
	queryDelete := `DELETE FROM likes WHERE likes_id = $1;`

	var boolValue bool
	var value, postId int64
	if err := repo.store.QueryRow(querySelect, likeId).Scan(&postId, &boolValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, repository.NotFound
		}
		return app.InvalidInt, repository.NewDBError(err)
	}

	if boolValue {
		value = 1
	} else {
		value = -1
	}

	begin, err := repo.store.Begin()
	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	var row *sql.Rows
	var newCountLikes int64 = -1
	err = begin.QueryRow(queryUpdate, postId, value).Scan(&newCountLikes)

	if err != nil {
		_ = begin.Rollback()
		return app.InvalidInt, repository.NewDBError(err)
	}

	if row, err = begin.Query(queryDelete, likeId); err != nil {
		_ = begin.Rollback()
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err = begin.Commit(); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return newCountLikes, nil
}
