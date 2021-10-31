package usercase_user

import (
	models2 "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type SuiteUserUsecase struct {
	usecase.SuiteUsecase
	uc    Usecase
	tUser *models.User
}

func (s *SuiteUserUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewUserUsecase(s.MockUserRepository, s.MockFilesRepository)
}

func (s *SuiteUserUsecase) TestCreatorUsecase_GetProfile_DB_Error() {
	s.Tb = usecase.TestTable{
		Name:              "DB error happened",
		Data:              int64(1),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.DefaultErrDB,
	}
	s.MockUserRepository.EXPECT().
		FindByID(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, repository.DefaultErrDB)
	u, err := s.uc.GetProfile(s.Tb.Data.(int64))
	assert.Nil(s.T(), u)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}

func (s *SuiteUserUsecase) TestCreatorUsecase_GetProfile_NotFound() {
	s.Tb = usecase.TestTable{
		Name:              "Profile not found",
		Data:              int64(1),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.NotFound,
	}
	s.MockUserRepository.EXPECT().
		FindByID(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, repository.NotFound)
	u, err := s.uc.GetProfile(s.Tb.Data.(int64))
	assert.Nil(s.T(), u)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_GetProfile_UserFound() {
	s.Tb = usecase.TestTable{
		Name:              "Profile found",
		Data:              int64(1),
		ExpectedMockTimes: 1,
		ExpectedError:     nil,
	}
	user := models.TestUser()
	s.MockUserRepository.EXPECT().
		FindByID(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(user, nil)
	u, err := s.uc.GetProfile(s.Tb.Data.(int64))
	assert.Equal(s.T(), user, u)
	assert.Equal(s.T(), s.Tb.ExpectedError, err)
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Check_NotFound() {
	s.Tb = usecase.TestTable{
		Name:              "User not found",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.NotFound,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "doggy123"
	req := models2.RequestLogin{
		Login:    u.Login,
		Password: u.Password,
	}
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, repository.NotFound)
	expectedId := int64(-1)
	resId, err := s.uc.Check(req.Login, req.Password)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Check_DB_Error() {
	s.Tb = usecase.TestTable{
		Name:              "Database error on request",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.DefaultErrDB,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "doggy123"
	req := models2.RequestLogin{
		Login:    u.Login,
		Password: u.Password,
	}
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, repository.DefaultErrDB)
	expectedId := int64(-1)
	resId, err := s.uc.Check(req.Login, req.Password)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Check_InvalidPassword() {
	s.Tb = usecase.TestTable{
		Name:              "Invalid password",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     models.IncorrectEmailOrPassword,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "doggy123"
	req := models2.RequestLogin{
		Login:    u.Login,
		Password: u.Password,
	}
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(u, nil)
	expectedId := int64(-1)
	resId, err := s.uc.Check(req.Login, req.Password)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}

func (s *SuiteUserUsecase) TestCreatorUsecase_Check_Correct() {
	s.Tb = usecase.TestTable{
		Name:              "User found, password valid",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     nil,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "doggy123"
	req := models2.RequestLogin{
		Login:    u.Login,
		Password: u.Password,
	}
	assert.NoError(s.T(), u.Encrypt())
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(u, nil)
	expectedId := u.ID
	resId, err := s.uc.Check(req.Login, req.Password)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, err)

}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_UserAlreadyExist() {
	s.Tb = usecase.TestTable{
		Name:              "User already exist",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     UserExist,
	}
	u := s.Tb.Data.(*models.User)
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(u, nil)
	expectedId := int64(-1)
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, err)
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_DB_Error() {
	s.Tb = usecase.TestTable{
		Name:              "Database error on find user",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.DefaultErrDB,
	}
	u := s.Tb.Data.(*models.User)
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, repository.DefaultErrDB)
	expectedId := int64(-1)
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_InvalidLoginShort() {
	s.Tb = usecase.TestTable{
		Name:              "Invalid login - short",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     models.IncorrectEmailOrPassword,
	}
	u := s.Tb.Data.(*models.User)
	u.Login = "l"
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	expectedId := int64(-1)
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_InvalidLoginLong() {
	s.Tb = usecase.TestTable{
		Name:              "Invalid login - long",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     models.IncorrectEmailOrPassword,
	}
	u := s.Tb.Data.(*models.User)
	u.Login = "llllllllllllllllllllllllllllllllllllll"
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	expectedId := int64(-1)
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_InvalidPasswordShort() {
	s.Tb = usecase.TestTable{
		Name:              "Invalid password - short",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     models.IncorrectEmailOrPassword,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "l"
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	expectedId := int64(-1)
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_InvalidPasswordLong() {
	s.Tb = usecase.TestTable{
		Name:              "Invalid password - long",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     models.IncorrectEmailOrPassword,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "lllllllllllllllllllllllllll" +
		"lllllllllllllllllllllllllll"
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	expectedId := int64(-1)
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_CreateFail() {
	s.Tb = usecase.TestTable{
		Name:              "Database error on create user",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.DefaultErrDB,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "lllllllllllllllllllllllllll"
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	s.MockUserRepository.EXPECT().
		Create(u).
		Times(s.Tb.ExpectedMockTimes).
		Return(repository.DefaultErrDB)
	expectedId := int64(-1)
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteUserUsecase) TestCreatorUsecase_Create_Success() {
	s.Tb = usecase.TestTable{
		Name:              "Success create user",
		Data:              models.TestUser(),
		ExpectedMockTimes: 1,
		ExpectedError:     nil,
	}
	u := s.Tb.Data.(*models.User)
	u.Password = "lllllllllllllllllllllllllll"
	s.MockUserRepository.EXPECT().
		FindByLogin(u.Login).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	s.MockUserRepository.EXPECT().
		Create(u).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil)
	expectedId := u.ID
	resId, err := s.uc.Create(u)
	assert.Equal(s.T(), expectedId, resId)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}

func TestUsecaseUser(t *testing.T) {
	suite.Run(t, new(SuiteUserUsecase))
}
