package usecase

import (
	"io"

	"github.com/ikmski/git-lfs3/entity"
)

type transferService struct {
	ContentRepository ContentRepository
	ContentPresenter  ContentPresenter
}

// TransferService is ...
type TransferService interface {
	Download(meta *entity.MetaData, w io.Writer) (int64, error)
	Upload(meta *entity.MetaData, r io.Reader) error
	Exists(meta *entity.MetaData) bool
}

// NewTransferService is ...
func NewTransferService(repo ContentRepository, pre ContentPresenter) TransferService {
	return &transferService{
		ContentRepository: repo,
		ContentPresenter:  pre,
	}
}

func (s *transferService) Download(meta *entity.MetaData, w io.Writer) (int64, error) {

	return 0, nil
}

func (s *transferService) Upload(meta *entity.MetaData, r io.Reader) error {

	return nil
}

func (s *transferService) Exists(meta *entity.MetaData) bool {

	return false
}
