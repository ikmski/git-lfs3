package usecase

import (
	"github.com/ikmski/git-lfs3/entity"
)

// BatchService is ...
type BatchService interface {
	Batch(req *BatchRequest) (*BatchResult, error)
}

type batchService struct {
	MetaDataRepository MetaDataRepository
	ContentRepository  ContentRepository
}

// NewBatchService is ...
func NewBatchService(metaDataRepo MetaDataRepository, contentRepo ContentRepository) BatchService {
	return &batchService{
		MetaDataRepository: metaDataRepo,
		ContentRepository:  contentRepo,
	}
}

func (c *batchService) Batch(req *BatchRequest) (*BatchResult, error) {

	var objectResults []*ObjectResult

	for _, obj := range req.Objects {

		meta, err := c.MetaDataRepository.Get(obj.Oid)

		if err == nil && c.ContentRepository.Exists(meta) {
			// Object is found and exists
			objectResult := createObjectResult(obj, meta, true, true)
			objectResults = append(objectResults, objectResult)
			continue
		}

		// Object is not found
		meta, err = c.MetaDataRepository.Put(obj.Oid, obj.Size)
		if err == nil {
			objectResult := createObjectResult(obj, meta, true, false)
			objectResults = append(objectResults, objectResult)
		}
	}

	result := &BatchResult{
		Objects: objectResults,
	}

	return result, nil
}

func createObjectResult(o *ObjectRequest, meta *entity.MetaData, metaExists, objectExists bool) *ObjectResult {

	return &ObjectResult{
		Oid:          meta.Oid,
		Size:         meta.Size,
		MetaExists:   metaExists,
		ObjectExists: objectExists,
	}
}
