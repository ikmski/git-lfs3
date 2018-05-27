package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

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

func (a *App) downloadHandler(c *gin.Context) {

	o := unpack(c)
	meta, err := a.metaStore.Get(o)
	if err != nil {
		writeStatus(c, 404)
		return
	}

	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {

		var fromByte int64 = 0
		var toByte int64 = meta.Size
		regex := regexp.MustCompile(`bytes=(.*)\-(.*)`)
		match := regex.FindStringSubmatch(rangeHeader)
		if match != nil && len(match) >= 3 {
			if len(match[1]) > 0 {
				fromByte, _ = strconv.ParseInt(match[1], 10, 64)
			}
			if len(match[2]) > 0 {
				toByte, _ = strconv.ParseInt(match[2], 10, 64)
			}
		}

		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", fromByte, toByte-1, int64(toByte-fromByte)))
		c.Status(206)
	}

	_, err = a.contentStore.Get(meta, newResponseWriterAt(c), rangeHeader)
	if err != nil {
		writeStatus(c, 404)
		return
	}
}

func (a *App) uploadHandler(c *gin.Context) {

	o := unpack(c)
	meta, err := a.metaStore.Get(o)
	if err != nil {
		writeStatus(c, 404)
		return
	}

	err = a.contentStore.Put(meta, c.Request.Body)
	if err != nil {
		a.metaStore.Delete(o)
		c.Status(500)
		fmt.Fprintf(c.Writer, `{"message":"%s"}`, err)
		return
	}
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

func unpack(c *gin.Context) *Object {

	o := &Object{
		User: c.Param("user"),
		Repo: c.Param("repo"),
		Oid:  c.Param("oid"),
	}

	return o
}

func writeStatus(c *gin.Context, status int) {

	message := http.StatusText(status)
	message = `{"message":"` + message + `"}`

	c.Status(status)
	fmt.Fprint(c.Writer, message)
}

type ResponseWriterAt struct {
	responseWriter http.ResponseWriter
}

func (rw *ResponseWriterAt) WriteAt(b []byte, off int64) (int, error) {
	return rw.responseWriter.Write(b)
}

func newResponseWriterAt(c *gin.Context) *ResponseWriterAt {
	return &ResponseWriterAt{
		responseWriter: c.Writer,
	}
}
