package adapter

import (
	"encoding/json"

	"github.com/ikmski/git-lfs3/usecase"
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

	req, err := parseBatchRequest(ctx)
	if err != nil {

	}

	result, err := c.BatchService.Batch(req)
	if err != nil {

	}

	res := convertBatchResponse(result)

	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.SetHeader("Content-Type", metaMediaType)
	ctx.SetStatus(200)
	ctx.GetResponseWriter().Write(json)
}

func parseBatchRequest(ctx Context) (*usecase.BatchRequest, error) {

	data, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}

	var req BatchRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	//user := ctx.GetParam("user")
	//repo := ctx.GetParam("repo")

	var objs []*usecase.ObjectRequest
	for _, o := range req.Objects {

		item := &usecase.ObjectRequest{
			Oid:  o.Oid,
			Size: o.Size,
		}

		objs = append(objs, item)
	}

	br := &usecase.BatchRequest{
		Objects: objs,
	}

	return br, nil
}

func convertBatchResponse(result *usecase.BatchResult) *BatchResponse {

	var objs []*ResponseObject

	for _, batchObj := range result.Objects {

		header := make(map[string]string)
		header["Accept"] = contentMediaType

		obj := newResponseObject()
		obj.Oid = batchObj.Oid
		obj.Size = batchObj.Size

		if batchObj.MetaExists {

		}

		if batchObj.ObjectExists {
			obj.Actions["download"] = &Link{
				Href:   "https://hoge",
				Header: header,
			}
		} else {
			obj.Actions["upload"] = &Link{
				Href:   "https://hoge",
				Header: header,
			}
		}

		objs = append(objs, obj)
	}

	res := &BatchResponse{
		Transfer: "basic",
		Objects:  objs,
	}

	return res
}
