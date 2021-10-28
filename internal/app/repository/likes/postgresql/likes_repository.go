package repository_postgresql

import (
	"database/sql"
	"github.com/pkg/errors"
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
	query := `SELECT posts_id, value, like_id FROM likes WHERE user_id = $1`
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

// Add Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *LikesRepository) Add(like *models.Like) error {
	query := `
		BEGIN;
		INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0);
		UPDATE posts SET posts.likes = posts.Likes + $3 WHERE posts_id = $1;
		COMMIT;`

	if _, err := repo.store.Query(query, like.PostId, like.UserId, like.Value); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// Delete Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *LikesRepository) Delete(likeId int64) error {
	querySelect := `SELECT posts_id, value FROM likes WHERE like_id = $1`
	queryDelete := `
		BEGIN;
		UPDATE posts SET posts.likes = posts.Likes - $1 WHERE posts_id = $2;
		DELETE FROM likes WHERE like_id = $3;
		COMMIT;`

	var value, postId int64
	if err := repo.store.QueryRow(querySelect, likeId).Scan(&postId, &value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	if _, err := repo.store.Query(queryDelete, value, postId, likeId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
