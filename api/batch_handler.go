package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikmski/git-lfs3/adapter"
)

type batchHandler struct {
	batchController adapter.BatchController
}

// BatchHandler is ...
type BatchHandler interface {
	Batch(c *gin.Context)
}

// NewBatchHandler is ...
func NewBatchHandler(c adapter.BatchController) BatchHandler {
	return &batchHandler{
		batchController: c,
	}
}

func (h *batchHandler) Batch(c *gin.Context) {

	br := parseBatchRequest(c)

	var responseObjects []*ResponseObject

	for _, object := range br.Objects {

		meta, err := a.metaStore.Get(object)

		if err == nil && a.contentStore.Exists(meta) {
			// Object is found and exists
			responseObject := a.createResponseObject(object, meta, true, false)
			responseObjects = append(responseObjects, responseObject)
			continue
		}

		// Object is not found
		meta, err = a.metaStore.Put(object)
		if err == nil {
			responseObject := a.createResponseObject(object, meta, meta.Existing, true)
			responseObjects = append(responseObjects, responseObject)
		}
	}

	response := &BatchResponse{
		Transfer: "basic",
		Objects:  responseObjects,
	}

	encoder := json.NewEncoder(c.Writer)
	encoder.Encode(response)

	c.Writer.Header().Set("Content-Type", metaMediaType)
	c.Status(200)
}

func createResponseObject(o *ObjectRequest, meta *ObjectMetaData, download, upload bool) *ResponseObject {

	rep := &ResponseObject{
		Oid:     meta.Oid,
		Size:    meta.Size,
		Actions: make(map[string]*Link),
	}

	header := make(map[string]string)
	header["Accept"] = contentMediaType

	if download {
		rep.Actions["download"] = &Link{
			Href:   o.DownloadLink(),
			Header: header,
		}
	}

	if upload {
		rep.Actions["upload"] = &Link{
			Href:   o.UploadLink(),
			Header: header,
		}
	}

	return rep
}

func randomLockId() string {

	var id [20]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}

func parseBatchRequest(c *gin.Context) *BatchRequest {

	var br BatchRequest

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&br)
	if err != nil {
		return &br
	}

	for i := 0; i < len(br.Objects); i++ {
		br.Objects[i].User = c.Param("user")
		br.Objects[i].Repo = c.Param("repo")
	}

	return &br
}

func writeStatus(c *gin.Context, status int) {

	message := http.StatusText(status)
	message = `{"message":"` + message + `"}`

	c.Status(status)
	fmt.Fprint(c.Writer, message)
}
