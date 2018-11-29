package usecase

import "time"

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

// ObjectError is ...
type ObjectError struct {
	Code    int
	Message string
}

// ObjectResult is ...
type ObjectResult struct {
	Oid          string
	Size         int64
	MetaExists   bool
	ObjectExists bool
	//Actions map[string]*Link
	//Error   *ObjectError
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
