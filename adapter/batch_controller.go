package adapter

import (
	"encoding/json"

	"github.com/ikmski/git-lfs3/usecase"
)

// BatchController is ...
type BatchController interface {
	Batch(ctx *Context) (*Context, error)
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

func (c *batchController) Batch(ctx *Context) {

	req := parseBatchRequest(ctx)

	res, err = c.batchService.Batch(req)

	if err != nil {

	}

	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.Header("Content-Type", metaMediaType)
	ctx.JSON(200, json)
}

func parseBatchRequest(ctx *Context) *usecase.BatchRequest {

	var br usecase.BatchRequest

	data, err := ctx.GetRawData()
	if err != nil {
		return &br
	}

	err = json.Unmarshal(data, &br)
	if err != nil {
		return &br
	}

	for i := 0; i < len(br.Objects); i++ {
		br.Objects[i].User = c.Param("user")
		br.Objects[i].Repo = c.Param("repo")
	}

	return &br
}
