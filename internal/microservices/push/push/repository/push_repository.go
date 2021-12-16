package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/repository"
)

const (
	GetUserNameAndAvatarQuery = `SELECT nickname, avatar FROM users WHERE users_id = $1`

	GetCreatorNameAndAvatarQuery = `SELECT u.nickname, cp.avatar FROM creator_profile as cp 
									JOIN users u on cp.creator_id = u.users_id
									WHERE cp.creator_id = $1`

	GetAwardsNameAndPriceQuery = `SELECT name, price FROM awards where awards_id = $1`

	GetCreatorPostAndTitleQuery = `SELECT creator_id, title FROM posts WHERE posts_id = $1`

	GetSubUserForPushPostQuery = `
					SELECT users_id FROM subscribers AS sb
					JOIN posts AS ps ON (ps.creator_id = sb.creator_id AND ps.posts_id = $1)
					JOIN user_settings AS us ON (us.user_id = users_id AND us.get_post)
					WHERE ps.is_draft = false AND (ps.type_awards is null OR ps.type_awards = sb.awards_id OR ps.type_awards IN
			 		(SELECT awa.awards_id FROM restapi_dev.public.parents_awards AS awa WHERE awa.parent_id = sb.awards_id))
	`
	CheckCreatorForGetSubPushQuery = `SELECT get_sub FROM user_settings WHERE user_id = $1`

	CheckCreatorForGetCommentPushQuery = `SELECT get_comment FROM user_settings WHERE user_id = $1`
)

type PushRepository struct {
	store *sqlx.DB
}

func NewPushRepository(st *sqlx.DB) *PushRepository {
	return &PushRepository{
		store: st,
	}
}

// GetUserNameAndAvatar Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetUserNameAndAvatar(userId int64) (nickname string, avatar string, err error) {
	if err = repo.store.QueryRow(GetUserNameAndAvatarQuery, userId).Scan(&nickname, &avatar); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", repository.NotFound
		}
		return "", "", repository.NewDBError(err)
	}

	return nickname, avatar, nil
}

// GetCreatorNameAndAvatar Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetCreatorNameAndAvatar(creatorId int64) (nickname string, avatar string, err error) {
	if err = repo.store.QueryRow(GetCreatorNameAndAvatarQuery, creatorId).Scan(&nickname, &avatar); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", repository.NotFound
		}
		return "", "", repository.NewDBError(err)
	}

	return nickname, avatar, nil
}

// GetAwardsNameAndPrice Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetAwardsNameAndPrice(awardsId int64) (name string, price int64, err error) {
	if err = repo.store.QueryRow(GetAwardsNameAndPriceQuery, awardsId).Scan(&name, &price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", app.InvalidInt, repository.NotFound
		}
		return "", app.InvalidInt, repository.NewDBError(err)
	}

	return name, price, nil
}

// GetCreatorPostAndTitle Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetCreatorPostAndTitle(postId int64) (int64, string, error) {
	var creatorId int64
	var title string
	if err := repo.store.QueryRow(GetCreatorPostAndTitleQuery, postId).Scan(&creatorId, &title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, "", repository.NotFound
		}
		return app.InvalidInt, "", repository.NewDBError(err)
	}

	return creatorId, title, nil
}

// GetSubUserForPushPost Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetSubUserForPushPost(postId int64) ([]int64, error) {
	var res []int64
	row, err := repo.store.Query(GetSubUserForPushPostQuery, postId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for row.Next() {
		var userId int64
		if err = row.Scan(&userId); err != nil {
			_ = row.Close()
			return nil, repository.NewDBError(err)
		}
		res = append(res, userId)
	}

	if err = row.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// CheckCreatorForGetSubPush Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) CheckCreatorForGetSubPush(creatorId int64) (bool, error) {
	var enable bool
	if err := repo.store.QueryRow(CheckCreatorForGetSubPushQuery, creatorId).Scan(&enable); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.NotFound
		}
		return false, repository.NewDBError(err)
	}

	return enable, nil
}

// CheckCreatorForGetCommentPush Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) CheckCreatorForGetCommentPush(creatorId int64) (bool, error) {
	var enable bool
	if err := repo.store.QueryRow(CheckCreatorForGetCommentPushQuery, creatorId).Scan(&enable); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.NotFound
		}
		return false, repository.NewDBError(err)
	}

	return enable, nil
}
