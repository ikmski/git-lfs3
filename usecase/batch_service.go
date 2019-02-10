package usecase

import "github.com/ikmski/git-lfs3/entity"

// BatchRequest is ...
type BatchRequest struct {
	Operation string
	Transfers []string
	Ref       string
	Objects   []*ObjectRequest
}

// ObjectRequest is ...
type ObjectRequest struct {
	Oid      string
	Size     int64
	User     string
	Password string
	Repo     string
}

// ObjectResult is ...
type ObjectResult struct {
	Oid          string
	Size         int64
	MetaExists   bool
	ObjectExists bool
}

// BatchResult is ...
type BatchResult struct {
	Transfer string
	Objects  []*ObjectResult
}

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

		meta, err := c.MetaDataRepository.Get(obj)

		if err == nil && c.ContentRepository.Exists(meta) {
			// Object is found and exists
			objectResult := createObjectResult(obj, meta, true, true)
			objectResults = append(objectResults, objectResult)
			continue
		}

		// Object is not found
		meta, err = c.MetaDataRepository.Put(obj)
		if err == nil {
			objectResult := createObjectResult(obj, meta, meta.Existing, false)
			objectResults = append(objectResults, objectResult)
		}
	}

	result := &BatchResult{
		Transfer: "basic",
		Objects:  objectResults,
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
