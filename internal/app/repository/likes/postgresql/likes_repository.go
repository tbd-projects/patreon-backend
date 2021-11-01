package repository_postgresql

import (
	"database/sql"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

type LikesRepository struct {
	store *sql.DB
}

func NewLikesRepository(st *sql.DB) *LikesRepository {
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
func (repo *LikesRepository) Add(like *models.Like) error {
	queryInsert := `INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0)`
	queryUpdate := `UPDATE posts SET likes = Likes + $2 WHERE posts_id = $1;`

	begin, err := repo.store.Begin()
	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}
	var row *sql.Rows
	row, err = repo.store.Query(queryInsert, like.PostId, like.UserId, like.Value)

	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	row, err = repo.store.Query(queryUpdate, like.PostId, like.Value)

	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = begin.Commit(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}
	return nil
}

// Delete Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *LikesRepository) Delete(likeId int64) error {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryDelete := `DELETE FROM likes WHERE likes_id = $1;`
	queryUpdate := `UPDATE posts SET likes = Likes - $2 WHERE posts_id = $1;`

	var boolValue bool
	var value, postId int64
	if err := repo.store.QueryRow(querySelect, likeId).Scan(&postId, &boolValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	if boolValue {
		value = 1
	} else {
		value = -1
	}

	begin, err := repo.store.Begin()
	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	var row *sql.Rows
	row, err = repo.store.Query(queryUpdate, postId, value)

	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if row, err = repo.store.Query(queryDelete, likeId); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = begin.Commit(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	return nil
}
