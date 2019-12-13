package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API is ...
type API struct {
    config *Config
	router *gin.Engine
}

type Config struct {
	Tls      bool
	Port     int
	Host     string
	CertFile string
	KeyFile  string
}

// NewAPI is ...
func NewAPI(
    conf *Config,
    batchHandler BatchHandler,
    transferHandler TransferHandler) *API {

	api := &API{
        config: conf,
    }

	r := gin.Default()

	r.POST("/:user/:repo/objects/batch", func (c *gin.Context) { batchHandler.Batch(c) })

	r.GET("/:user/:repo/objects/:oid", func (c *gin.Context) { transferHandler.Download(c) })
	r.PUT("/:user/:repo/objects/:oid", func (c *gin.Context) { transferHandler.Upload(c) })

	api.router = r

	return api
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	a.router.ServeHTTP(w, r)
}

// Serve is ...
func (a *API) Serve() error {

	if a.config.Tls {
		return a.router.RunTLS(
			fmt.Sprintf(":%d", a.config.Port),
			a.config.CertFile,
			a.config.KeyFile)
	}

	return a.router.Run(fmt.Sprintf(":%d", a.config.Port))
}
