package adapter

import "github.com/ikmski/git-lfs3/usecase"

type contentController struct {
	contentService usecase.ContentService
}

// ContentController is ...
type ContentController interface {
}

// NewContentController is ...
func NewContentController(s usecase.ContentService) ContentController {
	return &contentController{
		contentService: s,
	}
}
