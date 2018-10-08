package usecase

import (
	"io"

	"github.com/ikmski/git-lfs3/entity"
)

type contentService struct {
	ContentRepository ContentRepository
	ContentPresenter  ContentPresenter
}

// ContentService is ...
type ContentService interface {
	Download(meta *entity.MetaData, w io.Writer) (int64, error)
	Upload(meta *entity.MetaData, r io.Reader) error
	Exists(meta *entity.MetaData) bool
}

// NewContentService is ...
func NewContentService(repo ContentRepository, pre ContentPresenter) ContentService {
	return &contentService{
		ContentRepository: repo,
		ContentPresenter:  pre,
	}
}

func (s *contentService) Download(meta *entity.MetaData, w io.Writer) (int64, error) {

	return 0, nil
}

func (s *contentService) Upload(meta *entity.MetaData, r io.Reader) error {

	return nil
}

func (s *contentService) Exists(meta *entity.MetaData) bool {

	return false
}
