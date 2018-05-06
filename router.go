package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	router    *gin.Engine
	metaStore *MetaStore
}

func newApp(meta *MetaStore) *App {

	app := &App{metaStore: meta}

	r := gin.Default()

	r.POST("/{user}/{repo}/objects/batch", app.batchHandler)

	app.router = r

	return app
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	a.router.ServeHTTP(w, r)
}

func (a *App) Serve() error {

	if config.Server.Tls {
		return a.router.RunTLS(
			fmt.Sprintf(":%d", config.Server.Port),
			config.Server.CertFile,
			config.Server.KeyFile)
	} else {
		return a.router.Run(fmt.Sprintf(":%d", config.Server.Port))
	}
}
