package repository_postgresql

import (
	"database/sql"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	rp "patreon/internal/app/repository"
)


const (
	// Create
	queryCreate = `INSERT INTO creator_profile (creator_id, category,
		description, avatar, cover) VALUES ($1, $2, $3, $4, $5)
		RETURNING creator_id
	`
	queryCategoryCreate = `SELECT category_id FROM creator_category WHERE name = $1`

	// GetCreators
	queryCountGetCreators = `SELECT count(*) from creator_profile`
	queryCreatorGetCreators = `SELECT creator_id, cc.name, description, creator_profile.avatar, cover, usr.nickname 
					FROM creator_profile JOIN users AS usr ON usr.users_id = creator_profile.creator_id
					JOIN creator_category As cc ON creator_profile.category = cc.category_id`

	// GetCreator
	queryGetCreator = `SELECT cp.creator_id, cc.name, cp.description, cp.avatar, cp.cover, usr.nickname, sb.awards_id
			FROM creator_profile as cp JOIN users AS usr ON usr.users_id = cp.creator_id 
			JOIN creator_category As cc ON cp.category = cc.category_id 
			LEFT JOIN subscribers AS sb on (cp.creator_id = sb.creator_id and sb.users_id = $1)
			where cp.creator_id=$2`

	// ExistsCreator
	queryExistsCreator = `SELECT creator_id from creator_profile where creator_id=$1`

	// UpdateAvatar
	queryUpdateAvatar = `UPDATE creator_profile SET avatar = $1 WHERE creator_id = $2 RETURNING creator_id`

	// UpdateCover
	queryUpdateCover = `UPDATE creator_profile SET cover = $1 WHERE creator_id = $2 RETURNING creator_id`
)

type CreatorRepository struct {
	store *sql.DB
}

func NewCreatorRepository(st *sql.DB) *CreatorRepository {
	return &CreatorRepository{
		store: st,
	}
}

// Create Errors:
//		IncorrectCategory
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CreatorRepository) Create(cr *models.Creator) (int64, error) {
	category := int64(0)
	if err := repo.store.QueryRow(queryCategoryCreate, cr.Category).Scan(&category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, IncorrectCategory
		}
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err := repo.store.QueryRow(queryCreate, cr.ID, category, cr.Description,
		app.DefaultImage, app.DefaultImage).Scan(&cr.ID); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}
	return cr.ID, nil
}

// GetCreators Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) GetCreators() ([]models.Creator, error) {
	count := 0

	if err := repo.store.QueryRow(queryCountGetCreators).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	res := make([]models.Creator, count)

	rows, err := repo.store.Query(queryCreatorGetCreators)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var creator models.Creator
		if err = rows.Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}
		res[i] = creator
		i++
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// GetCreator Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) GetCreator(creatorId int64, userId int64) (*models.CreatorWithAwards, error) {
	creator := &models.CreatorWithAwards{}

	var awardsId sql.NullInt64
	if err := repo.store.QueryRow(queryGetCreator, userId, creatorId).
		Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname, &awardsId); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	if awardsId.Valid == false {
		creator.AwardsId = rp.NoAwards
	} else {
		creator.AwardsId = awardsId.Int64
	}

	return creator, nil
}

// ExistsCreator Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) ExistsCreator(creatorId int64) (bool, error) {
	creator := &models.Creator{}

	if err := repo.store.QueryRow(queryExistsCreator, creatorId).Scan(&creator.ID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, repository.NewDBError(err)
	}

	return true, nil
}

// UpdateAvatar Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) UpdateAvatar(creatorId int64, avatar string) error {
	if err := repo.store.QueryRow(queryUpdateAvatar, avatar, creatorId).
		Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// UpdateCover Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) UpdateCover(creatorId int64, cover string) error {
	if err := repo.store.QueryRow(queryUpdateCover, cover, creatorId).
		Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}
