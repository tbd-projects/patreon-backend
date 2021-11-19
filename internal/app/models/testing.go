package models

import (
	"errors"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/go-sqlxmock"
)

var BDError = errors.New("BD error")

type Suite struct {
	suite.Suite
	DB   *sqlx.DB
	Mock sqlmock.Sqlmock
}

func (s *Suite) InitBD() {
	s.T().Helper()

	var err error
	s.DB, s.Mock, err = sqlmock.Newx()
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestUser() *User {
	return &User{
		Login:    "student1999",
		Password: "1!2!3!",
		Nickname: "patron",
		Avatar:   "default.png",
	}

}
func TestUsers() []User {
	u1 := User{
		Login:    "test1",
		Nickname: "test1",
		Password: "123456",
	}
	u2 := User{
		Login:    "test2",
		Nickname: "test2",
		Password: "123456",
	}
	u3 := User{
		Login:    "test3",
		Nickname: "test3",
		Password: "123456",
	}
	return []User{u1, u2, u3}
}

func TestCreator() *Creator {
	return &Creator{
		ID:          1,
		Nickname:    "doggy2005",
		Category:    "podcasts",
		Description: "i love podcasts",
	}
}

func TestCreatorWithAwards() *CreatorWithAwards {
	return &CreatorWithAwards{
		ID:          1,
		Nickname:    "doggy2005",
		Category:    "podcasts",
		Description: "i love podcasts",
		AwardsId:    1,
	}
}

func TestCreatorSubscriber() *CreatorSubscribe {
	return &CreatorSubscribe{
		ID:          1,
		Nickname:    "doggy2005",
		Category:    "podcasts",
		Description: "i love podcasts",
		AwardsId:    1,
	}
}

func TestCreators() []Creator {
	cr1 := Creator{
		ID:          1,
		Category:    "podcasts",
		Description: "i love podcasts",
	}
	cr2 := Creator{
		ID:          2,
		Category:    "blog",
		Description: "i love podcasts",
	}
	cr3 := Creator{
		ID:          3,
		Category:    "movies",
		Description: "i love podcasts",
	}
	return []Creator{cr1, cr2, cr3}
}
func TestSubscriber() *Subscriber {
	return &Subscriber{
		ID:        1,
		UserID:    1,
		CreatorID: 2,
		AwardID:   3,
	}

}
func TestAward() *Award {
	return &Award{
		ID:          1,
		Name:        "award",
		Description: "description",
		Price:       100,
		CreatorId:   1,
		Cover:       "not found",
	}
}
func TestLike() *Like {
	return &Like{
		ID:     1,
		Value:  1,
		PostId: 2,
		UserId: 3,
	}
}
func TestUpdatePost() *UpdatePost {
	return &UpdatePost{
		ID:          1,
		Title:       "Title",
		Awards:      1,
		Description: "jfnagd",
	}
}
func TestCreatePost() *CreatePost {
	return &CreatePost{
		ID:          1,
		Title:       "Title",
		Awards:      1,
		Description: "jfnagd",
		CreatorId:   1,
	}
}
func TestAttachWithoutLevel() *AttachWithoutLevel {
	return &AttachWithoutLevel{
		ID:     1,
		PostId: 1,
		Value:   "jfnagd",
		Type:   Image,
	}
}
