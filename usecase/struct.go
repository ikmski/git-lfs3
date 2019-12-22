package usecase

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
