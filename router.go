package main

import "github.com/gin-gonic/gin"

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
