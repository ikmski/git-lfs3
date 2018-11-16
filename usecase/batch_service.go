package usecase

import "time"

// BatchService is ...
type BatchService interface {
	Batch(req *BatchRequest) (*BatchResult, error)
}

type batchService struct {
	MetaDataRepository MetaDataRepository
	MetaDataPresenter  MetaDataPresenter
}

// ObjectRequest is ...
type ObjectRequest struct {
	Oid      string
	Size     int64
	User     string
	Password string
	Repo     string
}

// BatchRequest is ...
type BatchRequest struct {
	Operation string
	Transfers []string
	Ref       string
	Objects   []*ObjectRequest
}

// ObjectError is ...
type ObjectError struct {
	Code    int
	Message string
}

// ObjectResult is ...
type ObjectResult struct {
	Oid     string
	Size    int64
	Actions map[string]*Link
	Error   *ObjectError
}

// BatchResult is ...
type BatchResult struct {
	Transfer string
	Objects  []*ObjectResult
}

// Link is ...
type Link struct {
	Href      string
	Header    map[string]string
	ExpiresAt time.Time
}

// NewBatchService is ...
func NewBatchService(repo MetaDataRepository, pre MetaDataPresenter) BatchService {
	return &batchService{
		MetaDataRepository: repo,
		MetaDataPresenter:  pre,
	}
}

func (c *batchService) Batch(req *BatchRequest) (*BatchResult, error) {

	return nil, nil
}
