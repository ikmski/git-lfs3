package usecase

import "github.com/ikmski/git-lfs3/entity"

// ContentPresenter is ...
type ContentPresenter interface {
	ResponseContent(c *entity.Content) *entity.Content
}
