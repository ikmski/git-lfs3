package adapter

import (
	"encoding/json"

	"github.com/ikmski/git-lfs3/usecase"
)

const (
	contentMediaType = "application/vnd.git-lfs"
	metaMediaType    = contentMediaType + "+json"
)

// BatchController is ...
type BatchController interface {
	Batch(ctx Context)
}

type batchController struct {
	BatchService usecase.BatchService
}

// NewBatchController is ...
func NewBatchController(s usecase.BatchService) BatchController {
	return &batchController{
		BatchService: s,
	}
}

func (c *batchController) Batch(ctx Context) {

	req := parseBatchRequest(ctx)

	result, err := c.BatchService.Batch(req)

	if err != nil {

	}

	res := convertBatchResponse(result)

	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.SetHeader("Content-Type", metaMediaType)
	ctx.SetJson(200, json)
}

func parseBatchRequest(ctx Context) *usecase.BatchRequest {

	var br BatchRequest

	data, err := ctx.GetRawData()
	if err != nil {
		return convertBatchRequest(&br)
	}

	err = json.Unmarshal(data, &br)
	if err != nil {
		return convertBatchRequest(&br)
	}

	for i := 0; i < len(br.Objects); i++ {
		br.Objects[i].User = ctx.GetParam("user")
		br.Objects[i].Repo = ctx.GetParam("repo")
	}

	return convertBatchRequest(&br)
}
