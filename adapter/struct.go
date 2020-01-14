package adapter

import (
	"time"
)

// BatchRequest is ...
type BatchRequest struct {
	Operation string           `json:"operation"`
	Transfers []string         `json:"transfers,omitempty"`
	Ref       Ref              `json:"ref,omitempty"`
	Objects   []*ObjectRequest `json:"objects"`
}

// ObjectRequest is ...
type ObjectRequest struct {
	Oid      string `json:"oid"`
	Size     int64  `json:"size"`
	User     string
	Password string
	Repo     string
}

// Ref is ...
type Ref struct {
	Name string `json:"name"`
}

// BatchResponse is ...
type BatchResponse struct {
	Transfer string            `json:"transfer,omitempty"`
	Objects  []*ResponseObject `json:"objects"`
}

// ResponseObject is ...
type ResponseObject struct {
	Oid     string           `json:"oid"`
	Size    int64            `json:"size"`
	Actions map[string]*Link `json:"actions"`
	Error   *ObjectError     `json:"error,omitempty"`
}

func newResponseObject() *ResponseObject {
	r := new(ResponseObject)
	r.Actions = make(map[string]*Link)
	return r
}

// Link is ...
type Link struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}

// ObjectError is ...
type ObjectError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
