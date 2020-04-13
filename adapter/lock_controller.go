package adapter

import "github.com/ikmski/git-lfs3/usecase"

// LockController is ...
type LockController interface {
	Lock(ctx Context)
	Unlock(ctx Context)
	List(ctx Context)
	Verify(ctx Context)
}

type lockController struct {
	LockService usecase.LockService
}

func NewLockController(s usecase.LockService) LockController {
	return &lockController{
		LockService: s,
	}
}

func (c *lockController) Lock(ctx Context) {

}

func (c *lockController) Unlock(ctx Context) {

}

func (c *lockController) List(ctx Context) {

}

func (c *lockController) Verify(ctx Context) {

}
