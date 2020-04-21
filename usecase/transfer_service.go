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
	GetSize(req *ObjectRequest) int64
}

// NewTransferService is ...
func NewTransferService(metaDataRepo MetaDataRepository, contentRepo ContentRepository) TransferService {
	return &transferService{
		ContentRepository:  contentRepo,
		MetaDataRepository: metaDataRepo,
	}
}

func (s *transferService) Download(req *ObjectRequest, w io.Writer) (int64, error) {

	meta, err := s.MetaDataRepository.Get(req.Oid)
	if err != nil {
		return 0, err
	}

	return s.ContentRepository.Get(meta, w, req.From, req.To)
}

func (s *transferService) Upload(req *ObjectRequest, r io.Reader) error {

	meta, err := s.MetaDataRepository.Get(req.Oid)
	if err != nil {
		return err
	}

	return s.ContentRepository.Put(meta, r)
}

func (s *transferService) Exists(req *ObjectRequest) bool {

	meta, err := s.MetaDataRepository.Get(req.Oid)
	if err != nil {
		return false
	}

	return s.ContentRepository.Exists(meta)
}

func (s *transferService) GetSize(req *ObjectRequest) int64 {

	meta, err := s.MetaDataRepository.Get(req.Oid)
	if err != nil {
		return 0
	}

	return meta.Size
}
