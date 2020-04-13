package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ikmski/git-lfs3/adapter"
)

type context struct {
	w http.ResponseWriter
	r *http.Request
}

func newContext(w http.ResponseWriter, r *http.Request) adapter.Context {
	ctx := new(context)
	ctx.w = w
	ctx.r = r
	return ctx
}

func (ctx *context) GetHeader(s string) string {
	return ctx.r.Header.Get(s)
}

func (ctx *context) GetParam(s string) string {
	vars := mux.Vars(ctx.r)
	return vars[s]
}

func (ctx *context) GetRawData() ([]byte, error) {
	buf := new(bytes.Buffer)
	io.Copy(buf, ctx.r.Body)
	return buf.Bytes(), nil
}

func (ctx *context) SetStatus(s int) {
	ctx.w.WriteHeader(s)
}

func (ctx *context) SetHeader(key string, val string) {
	ctx.w.Header().Set(key, val)
}

func (ctx *context) GetResponseWriter() io.Writer {
	return ctx.w
}

func (ctx *context) GetRequestReader() io.Reader {
	return ctx.r.Body
}
