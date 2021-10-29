package repository_os

type ImageRepository struct {
	staticDir string
}

func NewImageRepository(staticDir string) *ImageRepository {
	return &ImageRepository{
		staticDir: staticDir,
	}
}

