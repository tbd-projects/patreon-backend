package repository_files

import "io"

type TypeFiles string

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
