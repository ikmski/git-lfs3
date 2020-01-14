package usecase

// BatchRequest is ...
type BatchRequest struct {
	Objects []*ObjectRequest
}

// ObjectRequest is ...
type ObjectRequest struct {
	Oid  string
	Size int64
}

// BatchResult is ...
type BatchResult struct {
	Objects []*ObjectResult
}

// ObjectResult is ...
type ObjectResult struct {
	Oid          string
	Size         int64
	MetaExists   bool
	ObjectExists bool
}
