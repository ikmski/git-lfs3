package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (a *App) batchHandler(c *gin.Context) {

	br := unpackBatchRequest(c)

	var responseObjects []*Representation

	// Create a response object
	for _, object := range br.Objects {

		meta, err := a.metaStore.Get(object)
		if err == nil && a.contentStore.Exists(meta) {
			// Object is found and exists
			responseObjects = append(responseObjects, a.Represent(object, meta, true, false))
			continue
		}

		// Object is not found
		meta, err = a.metaStore.Put(object)
		if err == nil {
			responseObjects = append(responseObjects, a.Represent(object, meta, meta.Existing, true))
		}
	}

	c.Writer.Header().Set("Content-Type", metaMediaType)

	respobj := &BatchResponse{Objects: responseObjects}

	enc := json.NewEncoder(c.Writer)
	enc.Encode(respobj)
}

func (a *App) Represent(o *Object, meta *MetaObject, download, upload bool) *Representation {

	rep := &Representation{
		Oid:     meta.Oid,
		Size:    meta.Size,
		Actions: make(map[string]*link),
	}

	header := make(map[string]string)
	header["Accept"] = contentMediaType
	if download {
		rep.Actions["download"] = &link{Href: o.DownloadLink(), Header: header}
	}

	if upload {
		rep.Actions["upload"] = &link{Href: o.UploadLink(), Header: header}
	}
	return rep
}

func randomLockId() string {

	var id [20]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}

func unpackBatchRequest(c *gin.Context) *BatchRequest {

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
