package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikmski/git-lfs3/adapter"
)

type batchHandler struct {
	batchController adapter.BatchController
}

// BatchHandler is ...
type BatchHandler interface {
	Batch(c *gin.Context)
}

// NewBatchHandler is ...
func NewBatchHandler(c adapter.BatchController) BatchHandler {
	return &batchHandler{
		batchController: c,
	}
}

func (h *batchHandler) Batch(c *gin.Context) {

	h.batchController.Batch(c)
}

func writeStatus(c *gin.Context, status int) {

	message := http.StatusText(status)
	message = `{"message":"` + message + `"}`

	c.Status(status)
	fmt.Fprint(c.Writer, message)
}
