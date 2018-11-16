package adapter

import (
	"github.com/ikmski/git-lfs3/usecase"
)

// BatchController is ...
type BatchController interface {
	Batch(req *BatchRequest) (*BatchResult, error)
}

type batchController struct {
	batchService usecase.BatchService
}

// NewBatchController is ...
func NewBatchController(s usecase.BatchService) BatchController {
	return &batchController{
		batchService: s,
	}
}

func (c *batchController) Batch(req *BatchRequest) (*BatchResult, error) {

	return nil, nil
}
