package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	router       *gin.Engine
	metaStore    *MetaStore
	contentStore *ContentStore
}

func newApp(meta *MetaStore, content *ContentStore) *App {

	app := &App{
		metaStore:    meta,
		contentStore: content,
	}

	r := gin.Default()

	r.POST("/{user}/{repo}/objects/batch", app.batchHandler)

	r.GET("/{user}/{repo}/objects/{oid}", app.downloadHandler)

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
