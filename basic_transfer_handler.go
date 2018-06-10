package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *App) downloadHandler(c *gin.Context) {

	o := parseObjectRequest(c)
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

	o := parseObjectRequest(c)
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

func parseObjectRequest(c *gin.Context) *ObjectRequest {

	o := &ObjectRequest{
		User: c.Param("user"),
		Repo: c.Param("repo"),
		Oid:  c.Param("oid"),
	}

	return o
}
