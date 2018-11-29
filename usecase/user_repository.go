package usecase

import (
	"github.com/ikmski/git-lfs3/entity"
)

// UserRepository is ...
type UserRepository interface {
	AddUser(user, pass string) error
	DeleteUser(user string) error
	Users() ([]*entity.User, error)
}
