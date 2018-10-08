package usecase

import "github.com/ikmski/git-lfs3/entity"

// MetaDataPresenter is ...
type MetaDataPresenter interface {
	ResponseMetaData(m *entity.MetaData) *entity.MetaData
}
