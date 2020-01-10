package usecase

import (
	"github.com/ikmski/git-lfs3/entity"
)

// MetaDataRepository is ...
type MetaDataRepository interface {
	Get(o *ObjectRequest) (*entity.MetaData, error)
	Put(o *ObjectRequest) (*entity.MetaData, error)
	Delete(o *ObjectRequest) error
	Objects() ([]*entity.MetaData, error)
}
