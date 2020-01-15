package usecase

import (
	"io"

	"github.com/ikmski/git-lfs3/entity"
)

// ContentRepository is ...
type ContentRepository interface {
	Get(meta *entity.MetaData, w io.Writer, from int64, to int64) (int64, error)
	Put(meta *entity.MetaData, r io.Reader) error
	Exists(meta *entity.MetaData) bool
}
