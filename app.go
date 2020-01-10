package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikmski/git-lfs3/adapter"
)

type app struct {
	config serverConfig
	router *gin.Engine
}

func newApp(
	conf serverConfig,
	batchController adapter.BatchController,
	transferController adapter.TransferController) *app {

	a := &app{
		config: conf,
	}

	r := gin.Default()

	r.POST("/:user/:repo/objects/batch", func(c *gin.Context) { batchController.Batch(newContext(c)) })

	r.GET("/:user/:repo/objects/:oid", func(c *gin.Context) { transferController.Download(newContext(c)) })
	r.PUT("/:user/:repo/objects/:oid", func(c *gin.Context) { transferController.Upload(newContext(c)) })

	a.router = r

	return a
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	a.router.ServeHTTP(w, r)
}

func (a *app) serve() error {

	if a.config.Tls {
		return a.router.RunTLS(
			fmt.Sprintf(":%d", a.config.Port),
			a.config.CertFile,
			a.config.KeyFile)
	}

	return a.router.Run(fmt.Sprintf(":%d", a.config.Port))
}
