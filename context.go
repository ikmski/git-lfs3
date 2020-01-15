package main

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/ikmski/git-lfs3/adapter"
)

type context struct {
	c *gin.Context
}

func newContext(c *gin.Context) adapter.Context {
	ctx := new(context)
	ctx.c = c
	return ctx
}

func (ctx *context) GetHeader(s string) string {
	return ctx.c.GetHeader(s)
}

func (ctx *context) GetParam(s string) string {
	return ctx.c.Param(s)
}

func (ctx *context) GetRawData() ([]byte, error) {
	return ctx.c.GetRawData()
}

func (ctx *context) SetStatus(s int) {
	ctx.c.Status(s)
}

func (ctx *context) SetHeader(key string, val string) {
	ctx.c.Header(key, val)
}

func (ctx *context) GetResponseWriter() io.Writer {
	return ctx.c.Writer
}

func (ctx *context) GetRequestReader() io.Reader {
	return ctx.c.Request.Body
}
