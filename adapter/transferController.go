package adapter

import "github.com/ikmski/git-lfs3/usecase"

type transferController struct {
	transferService usecase.TransferService
}

// TransferController is ...
type TransferController interface {
}

// NewTransferController is ...
func NewTransferController(s usecase.TransferService) TransferController {
	return &transferController{
		transferService: s,
	}
}
