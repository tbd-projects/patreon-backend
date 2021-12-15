package attaches

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type SuiteAttachesUsecase struct {
	usecase.SuiteUsecase
	uc *AttachesUsecase

}

func (s *SuiteAttachesUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewAttachesUsecase(s.MockAttachesRepository, s.MockFileClient, s.MockConvector)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_GetAttach() {
	att := models.TestAttachWithoutLevel()

	s.MockAttachesRepository.EXPECT().
		Get(att.ID).
		Times(1).
		Return(att, nil)
	gotattach, err := s.uc.GetAttach(att.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), att, gotattach)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_Delete() {
	att := models.TestAttachWithoutLevel()

	s.MockAttachesRepository.EXPECT().
		Delete(att.ID).
		Times(1).
		Return(nil)
	err := s.uc.Delete(att.ID)
	assert.NoError(s.T(), err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_LoadImage() {
	att := models.TestAttachWithoutLevel()
	buff := bytes.NewBufferString("dor")
	reader := bytes.NewReader(buff.Bytes())
	fileName := repoFiles.FileName("dor")
	att.Value = app.LoadFileUrl + string(fileName)
	att.Type = models.Image
	att.ID = 0
	resId := int64(1)

	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Image).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, nil)
	id, err := s.uc.LoadImage(reader, fileName, att.PostId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resId, id)

	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, repository.DefaultErrDB)
	_, err = s.uc.LoadImage(reader, fileName, att.PostId)
	assert.Error(s.T(), err)

	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Image).
		Times(1).
		Return(string(fileName), repository.DefaultErrDB)
	_, err = s.uc.LoadImage(reader, fileName, att.PostId)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Image).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, repository.DefaultErrDB)
	_, err = s.uc.LoadImage(reader, fileName, att.PostId)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	att.PostId = -1
	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Image).
		Times(1).
		Return(string(fileName), nil)
	_, err = s.uc.LoadImage(reader, fileName, att.PostId)
	assert.Error(s.T(), err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_LoadVideo() {
	att := models.TestAttachWithoutLevel()
	buff := bytes.NewBufferString("dor")
	reader := bytes.NewReader(buff.Bytes())
	fileName := repoFiles.FileName("dor")
	att.Value = app.LoadFileUrl + string(fileName)
	att.Type = models.Video
	att.ID = 0
	resId := int64(1)

	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Video).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, nil)
	id, err := s.uc.LoadVideo(reader, fileName, att.PostId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resId, id)

	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Video).
		Times(1).
		Return(string(fileName), repository.DefaultErrDB)
	_, err = s.uc.LoadVideo(reader, fileName, att.PostId)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Video).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, repository.DefaultErrDB)
	_, err = s.uc.LoadVideo(reader, fileName, att.PostId)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	att.PostId = -1
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Video).
		Times(1).
		Return(string(fileName), nil)
	_, err = s.uc.LoadVideo(reader, fileName, att.PostId)
	assert.Error(s.T(), err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_LoadAudio() {
	att := models.TestAttachWithoutLevel()
	buff := bytes.NewBufferString("dor")
	reader := bytes.NewReader(buff.Bytes())
	fileName := repoFiles.FileName("dor")
	att.Value = app.LoadFileUrl + string(fileName)
	att.Type = models.Music
	att.ID = 0
	resId := int64(1)

	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Music).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, nil)
	id, err := s.uc.LoadAudio(reader, fileName, att.PostId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resId, id)

	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Music).
		Times(1).
		Return(string(fileName), repository.DefaultErrDB)
	_, err = s.uc.LoadAudio(reader, fileName, att.PostId)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Music).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, repository.DefaultErrDB)
	_, err = s.uc.LoadAudio(reader, fileName, att.PostId)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	att.PostId = -1
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Music).
		Times(1).
		Return(string(fileName), nil)
	_, err = s.uc.LoadAudio(reader, fileName, att.PostId)
	assert.Error(s.T(), err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_UpdateAudio() {
	att := models.TestAttachWithoutLevel()
	buff := bytes.NewBufferString("dor")
	reader := bytes.NewReader(buff.Bytes())
	fileName := repoFiles.FileName("dor")
	att.Value = app.LoadFileUrl + string(fileName)
	att.Type = models.Music
	att.ID = 1
	att.PostId = 0

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Music).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(nil)
	err := s.uc.UpdateAudio(reader, fileName, att.ID)
	assert.NoError(s.T(), err)

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Music).
		Times(1).
		Return(string(fileName), repository.DefaultErrDB)
	err = s.uc.UpdateAudio(reader, fileName, att.ID)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Music).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(repository.DefaultErrDB)
	err = s.uc.UpdateAudio(reader, fileName, att.ID)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, repository.DefaultErrDB)
	err = s.uc.UpdateAudio(reader, fileName, att.ID)
	assert.Error(s.T(), err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_UpdateVideo() {
	att := models.TestAttachWithoutLevel()
	buff := bytes.NewBufferString("dor")
	reader := bytes.NewReader(buff.Bytes())
	fileName := repoFiles.FileName("dor")
	att.Value = app.LoadFileUrl + string(fileName)
	att.Type = models.Video
	att.ID = 1
	att.PostId = 0

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Video).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(nil)
	err := s.uc.UpdateVideo(reader, fileName, att.ID)
	assert.NoError(s.T(), err)

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Video).
		Times(1).
		Return(string(fileName), repository.DefaultErrDB)
	err = s.uc.UpdateVideo(reader, fileName, att.ID)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Video).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(repository.DefaultErrDB)
	err = s.uc.UpdateVideo(reader, fileName, att.ID)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, repository.DefaultErrDB)
	err = s.uc.UpdateVideo(reader, fileName, att.ID)
	assert.Error(s.T(), err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_UpdateImage() {
	att := models.TestAttachWithoutLevel()
	buff := bytes.NewBufferString("dor")
	reader := bytes.NewReader(buff.Bytes())
	fileName := repoFiles.FileName("dor")
	att.Value = app.LoadFileUrl + string(fileName)
	att.Type = models.Image
	att.ID = 1
	att.PostId = 0

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Image).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(nil)
	err := s.uc.UpdateImage(reader, fileName, att.ID)
	assert.NoError(s.T(), err)

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Image).
		Times(1).
		Return(string(fileName), repository.DefaultErrDB)
	err = s.uc.UpdateImage(reader, fileName, att.ID)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, nil)
	s.MockFileClient.EXPECT().
		SaveFile(gomock.Any(), reader, fileName, repoFiles.Image).
		Times(1).
		Return(string(fileName), nil)
	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(repository.DefaultErrDB)
	err = s.uc.UpdateImage(reader, fileName, att.ID)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, repository.DefaultErrDB)
	err = s.uc.UpdateImage(reader, fileName, att.ID)
	assert.Error(s.T(), err)

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(att.ID).
		Times(1).
		Return(false, nil)
	s.MockConvector.EXPECT().
		Convert(gomock.Any(), reader, fileName).
		Times(1).
		Return(reader, fileName, repository.DefaultErrDB)
	err = s.uc.UpdateImage(reader, fileName, att.ID)
	assert.Error(s.T(), err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_LoadText() {
	att := models.TestAttachWithoutLevel()
	att.Type = models.Music
	att.ID = 0
	resId := int64(1)

	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, nil)
	id, err := s.uc.LoadText(att)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resId, id)

	s.MockAttachesRepository.EXPECT().
		Create(att).
		Times(1).
		Return(resId, repository.DefaultErrDB)
	_, err = s.uc.LoadText(att)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	att.PostId = -1
	_, err = s.uc.LoadText(att)
	assert.Error(s.T(), err)
}

func (s *SuiteAttachesUsecase) MockcheckAttach(updId int64) {
	s.MockAttachesRepository.EXPECT().
		ExistsAttach(updId).
		Times(1).
		Return(false, nil)
}

func (s *SuiteAttachesUsecase) MockcheckAttachError(updId int64, err error) {
	s.MockAttachesRepository.EXPECT().
		ExistsAttach(updId).
		Times(1).
		Return(false, err)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_checkAttach() {
	newAtt := []models.Attach{*models.TestAttach()}
	updAtt := []models.Attach{*models.TestAttach()}

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(updAtt[0].Id).
		Times(1).
		Return(false, nil)
	err := s.uc.checkAttach(newAtt, updAtt)
	assert.NoError(s.T(), err)

	s.MockAttachesRepository.EXPECT().
		ExistsAttach(updAtt[0].Id).
		Times(1).
		Return(false, nil)
	err = s.uc.checkAttach(newAtt, updAtt)
	assert.NoError(s.T(), err)

	newAtt[0].Id = -1
	err = s.uc.checkAttach(newAtt, updAtt)
	assert.ErrorIs(s.T(), err, models.IncorrectAttachId)

	newAtt[0].Id = 1
	updAtt[0].Id = -1
	err = s.uc.checkAttach(newAtt, updAtt)
	assert.ErrorIs(s.T(), err, models.IncorrectAttachId)

	updAtt[0].Id = 0
	err = s.uc.checkAttach(newAtt, updAtt)
	assert.ErrorIs(s.T(), err, models.IncorrectAttachId)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_UpdateAttaches() {
	newAtt := []models.Attach{*models.TestAttach()}
	updAtt := []models.Attach{*models.TestAttach()}
	postId := int64(3)
	res := []int64{1, 2}

	s.MockcheckAttach(updAtt[0].Id)
	s.MockAttachesRepository.EXPECT().
		ApplyChangeAttaches(postId, newAtt, updAtt).
		Times(1).
		Return(res, nil)
	got, err := s.uc.UpdateAttach(postId, newAtt, updAtt)
	assert.Equal(s.T(), res, got)
	assert.NoError(s.T(), err)

	s.MockcheckAttachError(updAtt[0].Id, repository.DefaultErrDB)
	_, err = s.uc.UpdateAttach(postId, newAtt, updAtt)
	assert.ErrorIs(s.T(), err, repository.DefaultErrDB)

	s.MockcheckAttach(updAtt[0].Id)
	s.MockAttachesRepository.EXPECT().
		ApplyChangeAttaches(postId, newAtt, updAtt).
		Times(1).
		Return(res, repository.DefaultErrDB)
	_, err = s.uc.UpdateAttach(postId, newAtt, updAtt)
	assert.ErrorIs(s.T(), err, repository.DefaultErrDB)
}

func (s *SuiteAttachesUsecase) TestCreatorUsecase_UpdateText() {
	att := models.TestAttachWithoutLevel()
	att.Type = models.Music

	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(nil)
	err := s.uc.UpdateText(att)
	assert.NoError(s.T(), err)

	s.MockAttachesRepository.EXPECT().
		Update(att).
		Times(1).
		Return(repository.DefaultErrDB)
	err = s.uc.UpdateText(att)
	assert.EqualError(s.T(), err, repository.DefaultErrDB.Error())

	att.PostId = -1
	err = s.uc.UpdateText(att)
	assert.Error(s.T(), err)
}

func TestUsecaseCreator(t *testing.T) {
	suite.Run(t, new(SuiteAttachesUsecase))
}
