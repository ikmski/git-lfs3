package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ikmski/git-lfs3/adapter"
)

type transferHandler struct {
	transferController adapter.TransferController
}

// TransferHandler is ...
type TransferHandler interface {
	Download(c *gin.Context)
	Upload(c *gin.Context)
}

// NewTransferHandler is ...
func NewTransferHandler(c adapter.TransferController) TransferHandler {
	return &transferHandler{
		transferController: c,
	}
}

// Download is ...
func (h *transferHandler) Download(c *gin.Context) {

	h.transferController.Download(c)
}

// Upload is ...
func (h *transferHandler) Upload(c *gin.Context) {

	h.transferController.Upload(c)
}

