package usecase_likes

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase"
	"testing"
)

type SuiteLikeUsecase struct {
	usecase.SuiteUsecase
	uc    Usecase
	tLike *models.Like
}

func (s *SuiteLikeUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewLikesUsecase(s.MockLikesRepository)
	s.tLike = &models.Like{ID: 2, PostId: 3, UserId: 1, Value: 1}
}

func (s *SuiteLikeUsecase) TestLikeUsecase_Add_Ok() {
	s.Tb = usecase.TestTable{
		Name:              "DB error happened",
		Data:              int64(1),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.DefaultErrDB,
	}

	like := &models.Like{ID: 1, PostId: 2, Value: 1, UserId: 1}
	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(int64(1), repository.NotFound)
	s.MockLikesRepository.EXPECT().
		Add(like).
		Times(s.Tb.ExpectedMockTimes).
		Return(like.ID, nil)
	u, err := s.uc.Add(like)
	assert.Equal(s.T(), u, like.ID)
	assert.NoError(s.T(), err)

	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(int64(1), repository.DefaultErrDB)
	_, err = s.uc.Add(like)
	assert.Equal(s.T(), u, like.ID)
	assert.Error(s.T(), err, repository.DefaultErrDB)

	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(int64(1), nil)
	_, err = s.uc.Add(like)
	assert.Equal(s.T(), u, like.ID)
	assert.Error(s.T(), err, repository.DefaultErrDB)

	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(int64(1), repository.NotFound)
	s.MockLikesRepository.EXPECT().
		Add(like).
		Times(s.Tb.ExpectedMockTimes).
		Return(like.ID, repository.DefaultErrDB)
	u, err = s.uc.Add(like)
	assert.Equal(s.T(), u, like.ID)
	assert.Error(s.T(), err, repository.DefaultErrDB)
}

func (s *SuiteLikeUsecase) TestLikeUsecase_Delete_Ok() {
	s.Tb = usecase.TestTable{
		Name:              "DB error happened",
		Data:              int64(1),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.DefaultErrDB,
	}

	like := &models.Like{ID: 1, PostId: 2, Value: 1, UserId: 1}
	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(like.ID, nil)
	s.MockLikesRepository.EXPECT().
		Delete(like.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(like.ID, nil)
	u, err := s.uc.Delete(like.PostId, like.UserId)
	assert.Equal(s.T(), u, like.ID)
	assert.NoError(s.T(), err)

	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(int64(1), repository.DefaultErrDB)
	_, err = s.uc.Delete(like.PostId, like.UserId)
	assert.Equal(s.T(), u, like.ID)
	assert.Error(s.T(), err, repository.DefaultErrDB)

	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(int64(1), repository.NotFound)
	_, err = s.uc.Delete(like.PostId, like.UserId)
	assert.Error(s.T(), err, IncorrectDelLike)

	s.MockLikesRepository.EXPECT().
		GetLikeId(like.UserId, like.PostId).
		Times(s.Tb.ExpectedMockTimes).
		Return(like.ID, nil)
	s.MockLikesRepository.EXPECT().
		Delete(like.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(like.ID, repository.DefaultErrDB)
	u, err = s.uc.Delete(like.PostId, like.UserId)
	assert.Equal(s.T(), u, like.ID)
	assert.Error(s.T(), err, repository.DefaultErrDB)
}

func TestUsecaseLike(t *testing.T) {
	suite.Run(t, new(SuiteLikeUsecase))
}
