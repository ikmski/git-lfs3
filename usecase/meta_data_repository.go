package usecase

import (
	"github.com/ikmski/git-lfs3/entity"
)

// MetaDataRepository is ...
type MetaDataRepository interface {
	Get(oid string) (*entity.MetaData, error)
	Put(oid string, size int64) (*entity.MetaData, error)
	Delete(oid string) error
	Objects() ([]*entity.MetaData, error)
}
