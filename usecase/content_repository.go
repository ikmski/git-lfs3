package usecase

import "github.com/ikmski/git-lfs3/entity"

// ContentRepository is ...
type ContentRepository interface {
	Exists(meta *entity.MetaData) bool
}
