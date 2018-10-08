package usecase

import "github.com/ikmski/git-lfs3/entity"

type metaDataService struct {
	MetaDataRepository MetaDataRepository
	MetaDataPresenter  MetaDataPresenter
}

// MetaDataService is ...
type MetaDataService interface {
	Get(oid string) (*entity.MetaData, error)
	Put(meta *entity.MetaData) error
	Delete(oid string) error
}

func NewMetaDataService(repo MetaDataRepository, pre MetaDataPresenter) MetaDataService {
	return &metaDataService{
		MetaDataRepository: repo,
		MetaDataPresenter:  pre,
	}
}

func (s *metaDataService) Get(oid string) (*entity.MetaData, error) {

	return nil, nil
}

func (s *metaDataService) Put(meta *entity.MetaData) error {

	return nil
}

func (s *metaDataService) Delete(oid string) error {

	return nil
}
