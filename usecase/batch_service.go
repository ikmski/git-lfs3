package usecase

import "github.com/ikmski/git-lfs3/entity"

const (
	contentMediaType = "application/vnd.git-lfs"
)

// BatchService is ...
type BatchService interface {
	Batch(req *BatchRequest) (*BatchResult, error)
}

type batchService struct {
	MetaDataRepository MetaDataRepository
	MetaDataPresenter  MetaDataPresenter
	ContentRepository  ContentRepository
	ContentPresenter   ContentPresenter
}

// NewBatchService is ...
func NewBatchService(
	metaDataRepo MetaDataRepository,
	metaDataPre MetaDataPresenter,
	contentRepo ContentRepository,
	contentPre ContentPresenter) BatchService {
	return &batchService{
		MetaDataRepository: metaDataRepo,
		MetaDataPresenter:  metaDataPre,
		ContentRepository:  contentRepo,
		ContentPresenter:   contentPre,
	}
}

func (c *batchService) Batch(req *BatchRequest) (*BatchResult, error) {

	var objectResults []*ObjectResult

	for _, obj := range req.Objects {

		meta, err := c.MetaDataRepository.Get(obj)

		if err == nil && c.ContentRepository.Exists(meta) {
			// Object is found and exists
			responseObject := createResponseObject(obj, meta, true, false)
			objectResults = append(objectResults, responseObject)
			continue
		}

		// Object is not found
		meta, err = c.MetaDataRepository.Put(obj)
		if err == nil {
			responseObject := createResponseObject(obj, meta, meta.Existing, true)
			objectResults = append(objectResults, responseObject)
		}
	}

	result := &BatchResult{
		Transfer: "basic",
		Objects:  objectResults,
	}

	return result, nil
}

func createResponseObject(o *ObjectRequest, meta *entity.MetaData, download, upload bool) *ObjectResult {

	rep := &ObjectResult{
		Oid:     meta.Oid,
		Size:    meta.Size,
		Actions: make(map[string]*Link),
	}

	/*
		header := make(map[string]string)
		header["Accept"] = contentMediaType

		if download {
			rep.Actions["download"] = &Link{
				Href:   o.DownloadLink(),
				Header: header,
			}
		}

		if upload {
			rep.Actions["upload"] = &Link{
				Href:   o.UploadLink(),
				Header: header,
			}
		}
	*/

	return rep
}
