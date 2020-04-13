package usecase

import (
	"github.com/ikmski/git-lfs3/entity"
)

// LockRepository is ...
type LockRepository interface {
	Add(repo string, l ...entity.Lock) error
	Delete(repo string, user string, id string, force bool) (*entity.Lock, error)
	Fetch(repo string) ([]entity.Lock, error)
	FilteredFetch(repo string, path string, cursor string, limit string) (locks []entity.Lock, next string, err error)
	FetchAll() ([]entity.Lock, error)
}
