package repository_files

import "io"

type TypeFiles string

//go:generate mockgen -destination=mocks/mock_files_repository.go -package=mock_repository -mock_names=Repository=FilesRepository . Repository

const (
	Image = TypeFiles("image")
	File  = TypeFiles("file")
	Video = TypeFiles("music")
	Music = TypeFiles("video")
)

type FileName string

type Repository interface {
	// SaveFile Errors:
	//		app.GeneralError Errors:
	//			repository_os.ErrorCreate
	//			repository_os.ErrorCopyFile
	SaveFile(file io.Reader, name FileName, typeF TypeFiles) (string, error)
}
