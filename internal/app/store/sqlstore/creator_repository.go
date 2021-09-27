package sqlstore

import (
	"database/sql"
	"patreon/internal/app/store"
	"patreon/internal/models"

	log "github.com/sirupsen/logrus"
)

type CreatorRepository struct {
	store          *Store
	UserRepository UserRepository
}

func NewCreatorRepository(st *Store) *CreatorRepository {
	return &CreatorRepository{
		store:          st,
		UserRepository: UserRepository{st},
	}
}

func (repo *CreatorRepository) Create(cr *models.Creator) error {
	if err := repo.store.db.QueryRow("INSERT INTO creator_profile (creator_id, category, "+
		"description, avatar, cover) VALUES ($1, $2, $3, $4, $5)"+
		"RETURNING creator_id", cr.ID, cr.Category, cr.Description, cr.Avatar, cr.Cover).Scan(&cr.ID); err != nil {
		return err
	}
	return nil
}

func (repo *CreatorRepository) GetCreators() ([]models.Creator, error) {
	count := 0

	if err := repo.store.db.QueryRow("SELECT count(*) from creator_profile").Scan(&count); err != nil {
		return nil, err
	}
	res := make([]models.Creator, count)

	rows, err := repo.store.db.Query("SELECT creator_id, cover, description, category, avatar, usr.nickname " +
		"from creator_profile join user as usr on usr.user_id = creator_profile.creator_id")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	i := 0
	for empty := rows.Next(); !empty; i++ {
		var creator models.Creator
		if err = rows.Scan(&creator.ID, &creator.Cover, &creator.Description, &creator.Category,
			&creator.Avatar, &creator.Nickname); err != nil {
			return nil, err
		}
		res[i] = creator

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return res, err
}

func (repo *CreatorRepository) GetCreator(creatorId int64) (*models.Creator, error) {
	creator := &models.Creator{}

	if err := repo.store.db.QueryRow("SELECT creator_id, cover, description, category, background, usr.nickname " +
		"from creator_profile join user as usr on usr.user_id = creator_profile.creator_id where creator_id=$1", creatorId).
		Scan(&creator.ID, &creator.Cover, &creator.Description, &creator.Category,
			&creator.Avatar, &creator.Nickname); err != nil {
		return nil, store.NotFound
	}

	return creator, nil
}

//func (repo *CreatorRepository) FindByLogin(login string) (*models.User, error) {
//	user := models.User{}
//
//	if err := repo.store.db.QueryRow("SELECT user_id, login, encrypted_password from users where login=$1", login).
//		Scan(&user.ID, &user.Login, &user.EncryptedPassword); err != nil {
//		return nil, store.NotFound
//	}
//
//	return &user, nil
//}

//func (repo *CreatorRepository) FindByID(id int64) (*models.User, error) {
//	user := models.User{}
//
//	if err := repo.store.db.QueryRow("SELECT user_id, nickname, avatar from users where user_id=$1", id).
//		Scan(&user.ID, &user.Nickname, &user.Avatar); err != nil {
//		return nil, store.NotFound
//	}
//
//	return &user, nil
//}
