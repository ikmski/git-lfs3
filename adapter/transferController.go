package adapter

import "github.com/ikmski/git-lfs3/usecase"

type transferController struct {
	contentService usecase.ContentService
}

// TransferController is ...
type TransferController interface {
}

// NewTransferController is ...
func NewTransferController(s usecase.ContentService) TransferController {
	return &transferController{
		contentService: s,
	}
}
