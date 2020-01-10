package usecase

import (
	"github.com/ikmski/git-lfs3/entity"
)

// LockRepository is ...
type LockRepository interface {
	AddLocks(repo string, l ...entity.Lock) error
	Locks(repo string) ([]entity.Lock, error)
	FilteredLocks(repo, path, cursor, limit string) (locks []entity.Lock, next string, err error)
	DeleteLock(repo, user, id string, force bool) (*entity.Lock, error)
	AllLocks() ([]entity.Lock, error)
}
