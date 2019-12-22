package usecase

import (
	"io"
)

type transferService struct {
	ContentRepository  ContentRepository
	MetaDataRepository MetaDataRepository
}

// TransferService is ...
type TransferService interface {
	Download(req *ObjectRequest, w io.Writer) (int64, error)
	Upload(req *ObjectRequest, r io.Reader) error
	Exists(req *ObjectRequest) bool
}

// NewTransferService is ...
func NewTransferService(metaDataRepo MetaDataRepository, contentRepo ContentRepository) TransferService {
	return &transferService{
		ContentRepository:  contentRepo,
		MetaDataRepository: metaDataRepo,
	}
}

func (s *transferService) Download(req *ObjectRequest, w io.Writer) (int64, error) {

	return 0, nil
}

func (s *transferService) Upload(req *ObjectRequest, r io.Reader) error {

	return nil
}

func (s *transferService) Exists(req *ObjectRequest) bool {

	return false
}
