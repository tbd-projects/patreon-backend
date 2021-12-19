package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/repository"
	"strings"
)

const (
	getUserNameAndAvatarQuery = `SELECT nickname, avatar FROM users WHERE users_id = $1`

	getCreatorNameAndAvatarQuery = `SELECT u.nickname, cp.avatar FROM creator_profile as cp 
									JOIN users u on cp.creator_id = u.users_id
									WHERE cp.creator_id = $1`

	getAwardsNameAndPriceQuery = `SELECT name, price FROM awards where awards_id = $1`

	getCreatorPostAndTitleQuery = `SELECT creator_id, title FROM posts WHERE posts_id = $1`

	getSubUserForPushPostQuery = `
					SELECT users_id FROM subscribers AS sb
					JOIN posts AS ps ON (ps.creator_id = sb.creator_id AND ps.posts_id = $1)
					JOIN user_settings AS us ON (us.user_id = users_id AND us.get_post)
					WHERE ps.is_draft = false AND (ps.type_awards is null OR ps.type_awards = sb.awards_id OR ps.type_awards IN
			 		(SELECT awa.awards_id FROM restapi_dev.public.parents_awards AS awa WHERE awa.parent_id = sb.awards_id))
	`
	checkCreatorForGetSubPushQuery = `SELECT get_sub FROM user_settings WHERE user_id = $1`

	checkCreatorForGetCommentPushQuery = `SELECT get_comment FROM user_settings WHERE user_id = $1`

	getAwardsInfoAndCreatorIdAndUserIdFromPaymentsQuery = `SELECT p.creator_id, p.awards_id, a.name, p.users_id FROM payments as p 
												 JOIN awards a on p.awards_id = a.awards_id WHERE pay_token = $1`

	addPushInfoQuery = `INSERT INTO push_history (users_id, push_type, push) VALUES `

	getPushInfoQuery = `SELECT id, push_type, push, date, is_viewed FROM push_history WHERE users_id = $1`

	markViewedQuery = `UPDATE push_history SET is_viewed = true WHERE id = $1 and users_id = $2`
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
	if err = repo.store.QueryRow(getUserNameAndAvatarQuery, userId).Scan(&nickname, &avatar); err != nil {
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
	if err = repo.store.QueryRow(getCreatorNameAndAvatarQuery, creatorId).Scan(&nickname, &avatar); err != nil {
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
	if err = repo.store.QueryRow(getAwardsNameAndPriceQuery, awardsId).Scan(&name, &price); err != nil {
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
	if err := repo.store.QueryRow(getCreatorPostAndTitleQuery, postId).Scan(&creatorId, &title); err != nil {
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
	row, err := repo.store.Query(getSubUserForPushPostQuery, postId)
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
	if err := repo.store.QueryRow(checkCreatorForGetSubPushQuery, creatorId).Scan(&enable); err != nil {
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
	if err := repo.store.QueryRow(checkCreatorForGetCommentPushQuery, creatorId).Scan(&enable); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.NotFound
		}
		return false, repository.NewDBError(err)
	}

	return enable, nil
}

// GetAwardsInfoAndCreatorIdAndUserIdFromPayments Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetAwardsInfoAndCreatorIdAndUserIdFromPayments(token string) (*PaymentsInfo, error) {
	res := &PaymentsInfo{}
	if err := repo.store.QueryRow(getAwardsInfoAndCreatorIdAndUserIdFromPaymentsQuery, token).
		Scan(&res.CreatorId, &res.AwardsId, &res.AwardsName, &res.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// AddPushInfo Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) AddPushInfo(userId []int64, pushType string, push interface{}) error {
	var argsString []string
	var args []interface{}
	bdIndex := 1
	for _, id := range userId {
		ph := PushJson{Push: push}
		res, err := easyjson.Marshal(ph)
		if err != nil {
			continue
		}

		argsString = append(argsString, "(?, ?, ?)")
		args = append(args, id)
		args = append(args, pushType)
		args = append(args, types.JSONText(res))

		bdIndex += 4
	}

	query := fmt.Sprintf("%s %s", addPushInfoQuery,
		strings.Join(argsString, ", "))
	query = repo.store.Rebind(query)

	if _, err := repo.store.Exec(query, args...); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// GetPushInfo Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetPushInfo(userId int64) ([]Push, error) {
	var res []Push
	row, err := repo.store.Query(getPushInfoQuery, userId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for row.Next() {
		var currentPush Push
		var json types.JSONText
		if err = row.Scan(&currentPush.Id, &currentPush.Type, &json, &currentPush.Date, &currentPush.Viewed); err != nil {
			_ = row.Close()
			return nil, repository.NewDBError(err)
		}
		err = easyjson.Unmarshal(json, &currentPush.Push)
		if err != nil {
			continue
		}
		res = append(res, currentPush)
	}

	if err = row.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// MarkViewed Errors:
//		repository.NotModify
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) MarkViewed(pushId int64, userId int64) error {
	if res, err := repo.store.Exec(markViewedQuery, pushId, userId); err != nil {
		return repository.NewDBError(err)
	} else {
		if affected, err := res.RowsAffected(); err != nil {
			return repository.NewDBError(err)
		} else {
			if affected == 0 {
				return NotModify
			}
		}
	}
	return nil
}