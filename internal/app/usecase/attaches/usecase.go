package attaches

import (
	"io"
	"patreon/internal/app/models"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
)

const (
	BaseLimit = 10
	EmptyUser = -2
)

//go:generate mockgen -destination=mocks/mock_attaches_usecase.go -package=mock_usecase -mock_names=Usecase=AttachesUsecase . Usecase

type Usecase interface {

	// GetAttach Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAttach(attachId int64) (*models.AttachWithoutLevel, error)

	// UpdateAttach Errors:
	//		repository.NotFound
	//		repository_postgresql.UnknownDataFormat
	//		models.IncorrectType
	//  	models.IncorrectAttachId
	//      models.IncorrectLevel
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	UpdateAttach(postId int64, newAttach []models.Attach, updatedAttach []models.Attach) ([]int64, error)

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(postId int64) error

	// LoadImage Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	LoadImage(data io.Reader, name repoFiles.FileName, postId int64) (int64, error)

	// LoadAudio Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	LoadAudio(data io.Reader, name repoFiles.FileName, postId int64) (int64, error)

	// LoadVideo Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	LoadVideo(data io.Reader, name repoFiles.FileName, postId int64) (int64, error)

	// LoadText Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	LoadText(postData *models.AttachWithoutLevel) (int64, error)

	// UpdateText Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository.NotFound
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	UpdateText(postData *models.AttachWithoutLevel) error

	// UpdateImage Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository.NotFound
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	UpdateImage(data io.Reader, name repoFiles.FileName, AttachId int64) error

	// UpdateAudio Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository.NotFound
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	UpdateAudio(data io.Reader, name repoFiles.FileName, AttachId int64) error

	// UpdateVideo Errors:
	//		models.InvalidPostId
	//		models.InvalidType
	//		repository.NotFound
	//		repository_postgresql.UnknownDataFormat
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	UpdateVideo(data io.Reader, name repoFiles.FileName, AttachId int64) error
}
