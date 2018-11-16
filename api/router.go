package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API is ...
type API struct {
	router *gin.Engine
}

// NewAPI is ...
func NewAPI() *API {

	api := &App

	r := gin.Default()

	r.POST("/:user/:repo/objects/batch", batchHandler)

	r.GET("/:user/:repo/objects/:oid", downloadHandler)
	r.PUT("/:user/:repo/objects/:oid", uploadHandler)

	api.router = r

	return api
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	a.router.ServeHTTP(w, r)
}

// Serve is ...
func (a *API) Serve() error {

	if config.Server.Tls {
		return a.router.RunTLS(
			fmt.Sprintf(":%d", config.Server.Port),
			config.Server.CertFile,
			config.Server.KeyFile)
	}

	return a.router.Run(fmt.Sprintf(":%d", config.Server.Port))
}
